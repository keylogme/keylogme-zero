#include "keylogger_darwin.h"
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/hid/IOHIDKeys.h>
#include <IOKit/hid/IOHIDManager.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
  IOHIDDeviceRef device;
  int vendorID;
  int productID;
} HIDDeviceContext;

#define MAX_DEVICES 64
static IOHIDManagerRef hidManager = NULL;
static HIDDeviceContext *contexts[MAX_DEVICES];
static int contextCount = 0;
static bool runLoopStarted = false;

void HIDCallback(void *context, IOReturn result, void *sender,
                 IOHIDValueRef value) {
  HIDDeviceContext *ctx = (HIDDeviceContext *)context;
  IOHIDElementRef element = IOHIDValueGetElement(value);
  uint32_t usagePage = IOHIDElementGetUsagePage(element);
  uint32_t usage = IOHIDElementGetUsage(element);
  CFIndex pressed = IOHIDValueGetIntegerValue(value);

  if (usagePage == kHIDPage_KeyboardOrKeypad && usage >= 4 && usage <= 231) {
    printf("Key usage: 0x%x, pressed: %ld\n", usage, pressed);

    GoHandleKeyEvent(usage, (int)pressed, ctx->vendorID, ctx->productID);
  } else if (usagePage == kHIDPage_Button && usage >= 1 && usage <= 8) {
    printf("Click usage: 0x%x, pressed: %ld\n", usage, pressed);
    GoHandleKeyEvent(usage, (int)pressed, ctx->vendorID, ctx->productID);
  }
}

void DeviceRemovalCallback(void *context, IOReturn result, void *sender) {
  HIDDeviceContext *ctx = (HIDDeviceContext *)context;
  if (!ctx)
    return;

  IOHIDDeviceUnscheduleFromRunLoop(ctx->device, CFRunLoopGetCurrent(),
                                   kCFRunLoopDefaultMode);

  GoHandleDeviceEvent(ctx->vendorID, ctx->productID, 0); // disconnected

  for (int i = 0; i < contextCount; i++) {
    if (contexts[i] == ctx) {
      free(contexts[i]);
      contexts[i] = NULL;
      break;
    }
  }
}

void ManagerDeviceRemovalCallback(void *context, IOReturn result, void *sender,
                                  IOHIDDeviceRef device) {
  DeviceRemovalCallback(context, result, sender);
}

void DeviceMatchingCallback(void *context, IOReturn result, void *sender,
                            IOHIDDeviceRef device) {
  if (contextCount >= MAX_DEVICES)
    return;

  CFNumberRef vendorRef =
      IOHIDDeviceGetProperty(device, CFSTR(kIOHIDVendorIDKey));
  CFNumberRef productRef =
      IOHIDDeviceGetProperty(device, CFSTR(kIOHIDProductIDKey));
  int v = 0, p = 0;
  if (vendorRef)
    CFNumberGetValue(vendorRef, kCFNumberIntType, &v);
  if (productRef)
    CFNumberGetValue(productRef, kCFNumberIntType, &p);

  for (int i = 0; i < contextCount; i++) {
    if (contexts[i] && contexts[i]->device == device)
      return; // already added
  }

  HIDDeviceContext *ctx = malloc(sizeof(HIDDeviceContext));
  ctx->device = device;
  ctx->vendorID = v;
  ctx->productID = p;

  IOHIDDeviceRegisterInputValueCallback(device, HIDCallback, ctx);
  IOHIDDeviceScheduleWithRunLoop(device, CFRunLoopGetCurrent(),
                                 kCFRunLoopDefaultMode);
  IOHIDDeviceRegisterRemovalCallback(device, DeviceRemovalCallback, ctx);

  contexts[contextCount++] = ctx;

  GoHandleDeviceEvent(v, p, 1); // connected
}

CFMutableDictionaryRef CreateMatchingDict(int vendorID, int productID) {
  CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
      kCFAllocatorDefault, 0, &kCFTypeDictionaryKeyCallBacks,
      &kCFTypeDictionaryValueCallBacks);

  CFNumberRef vendorNum =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &vendorID);
  CFNumberRef productNum =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &productID);

  CFDictionarySetValue(dict, CFSTR(kIOHIDVendorIDKey), vendorNum);
  CFDictionarySetValue(dict, CFSTR(kIOHIDProductIDKey), productNum);

  CFRelease(vendorNum);
  CFRelease(productNum);
  return dict;
}

void setupDevice(int vendorID, int productID) {
  if (!hidManager) {
    hidManager = IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
    IOHIDManagerRegisterDeviceMatchingCallback(hidManager,
                                               DeviceMatchingCallback, NULL);
    IOHIDManagerRegisterDeviceRemovalCallback(
        hidManager, ManagerDeviceRemovalCallback, NULL);
    IOHIDManagerScheduleWithRunLoop(hidManager, CFRunLoopGetCurrent(),
                                    kCFRunLoopDefaultMode);
    IOHIDManagerOpen(hidManager, kIOHIDOptionsTypeNone);
  }

  CFDictionaryRef dict = CreateMatchingDict(vendorID, productID);
  IOHIDManagerSetDeviceMatchingMultiple(
      hidManager,
      CFArrayCreate(NULL, (const void **)&dict, 1, &kCFTypeArrayCallBacks));
  CFRelease(dict);
}

void Start(void) {
  if (hidManager && !runLoopStarted) {
    runLoopStarted = true;
    CFRunLoopRun();
    runLoopStarted = false;
  }
}

void Stop(void) {
  for (int i = 0; i < contextCount; i++) {
    if (contexts[i]) {
      IOHIDDeviceUnscheduleFromRunLoop(
          contexts[i]->device, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
      free(contexts[i]);
      contexts[i] = NULL;
    }
  }
  contextCount = 0;

  if (hidManager) {
    IOHIDManagerUnscheduleFromRunLoop(hidManager, CFRunLoopGetCurrent(),
                                      kCFRunLoopDefaultMode);
    IOHIDManagerClose(hidManager, kIOHIDOptionsTypeNone);
    CFRelease(hidManager);
    hidManager = NULL;
  }

  if (runLoopStarted) {
    CFRunLoopStop(CFRunLoopGetCurrent());
    runLoopStarted = false;
  }
}

void ListConnectedHIDDevices() {
  IOHIDManagerRef mgr =
      IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
  IOHIDManagerSetDeviceMatching(mgr, NULL); // Match all devices
  IOHIDManagerOpen(mgr, kIOHIDOptionsTypeNone);

  CFSetRef deviceSet = IOHIDManagerCopyDevices(mgr);
  if (!deviceSet) {
    printf("No HID devices found.\n");
    return;
  }

  CFIndex count = CFSetGetCount(deviceSet);
  IOHIDDeviceRef *devices = malloc(sizeof(IOHIDDeviceRef) * count);
  CFSetGetValues(deviceSet, (const void **)devices);

  // Use a simple string array to track unique keys
  char **seenKeys = calloc(count, sizeof(char *));
  size_t seenCount = 0;

  for (CFIndex i = 0; i < count; i++) {
    IOHIDDeviceRef dev = devices[i];

    CFStringRef nameRef = IOHIDDeviceGetProperty(dev, CFSTR(kIOHIDProductKey));
    CFNumberRef vendorRef =
        IOHIDDeviceGetProperty(dev, CFSTR(kIOHIDVendorIDKey));
    CFNumberRef productRef =
        IOHIDDeviceGetProperty(dev, CFSTR(kIOHIDProductIDKey));

    char name[256] = "(unknown)";
    if (nameRef) {
      CFStringGetCString(nameRef, name, sizeof(name), kCFStringEncodingUTF8);
    }

    int vid = 0, pid = 0;
    if (vendorRef)
      CFNumberGetValue(vendorRef, kCFNumberIntType, &vid);
    if (productRef)
      CFNumberGetValue(productRef, kCFNumberIntType, &pid);

    // Create a uniqueness key: "vid:pid:name"
    char key[512];
    snprintf(key, sizeof(key), "%04x:%04x:%s", vid, pid, name);

    // Check if this key has already been seen
    int duplicate = 0;
    for (size_t j = 0; j < seenCount; j++) {
      if (strcmp(seenKeys[j], key) == 0) {
        duplicate = 1;
        break;
      }
    }
    if (duplicate)
      continue;

    // Print and record new key
    printf("Device: %s | Vendor ID: 0x%04x | Product ID: 0x%04x\n", name, vid,
           pid);
    seenKeys[seenCount++] = strdup(key);
  }

  // Free memory
  for (size_t i = 0; i < seenCount; i++)
    free(seenKeys[i]);
  free(seenKeys);
  free(devices);
  CFRelease(deviceSet);
  CFRelease(mgr);
}
