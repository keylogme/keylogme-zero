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
  IOHIDManagerRef manager;
  CFRunLoopRef runLoop;
} HIDDeviceContext;

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

void stopDevice(CFRunLoopRef runLoop, IOHIDManagerRef manager) {

  IOHIDManagerUnscheduleFromRunLoop(manager, runLoop, kCFRunLoopDefaultMode);
  IOHIDManagerClose(manager, kIOHIDOptionsTypeNone);
  CFRelease(manager);

  // Stop loop
  CFRunLoopPerformBlock(runLoop, kCFRunLoopDefaultMode, ^{
    CFRunLoopStop(runLoop);
  });
  CFRunLoopWakeUp(runLoop);
  printf("Finished stopping device\n");
}

void DeviceRemovalCallback(void *context, IOReturn result, void *sender) {
  HIDDeviceContext *ctx = (HIDDeviceContext *)context;
  printf("Device removal callback\n");
  if (!ctx)
    return;

  IOHIDDeviceUnscheduleFromRunLoop(ctx->device, CFRunLoopGetCurrent(),
                                   kCFRunLoopDefaultMode);

  GoHandleDeviceEvent(ctx->vendorID, ctx->productID, 0); // disconnected

  // cleanup
  stopDevice(ctx->runLoop, ctx->manager);
}

void ManagerDeviceRemovalCallback(void *context, IOReturn result, void *sender,
                                  IOHIDDeviceRef device) {
  DeviceRemovalCallback(context, result, sender);
}

// Device matching callback may be called multiple times for the same device
// because a device can have multiple interfaces (keyboards, controls, mouse
// integrated, etc )
void DeviceMatchingCallback(void *context, IOReturn result, void *sender,
                            IOHIDDeviceRef device) {
  printf("Device matching callback\n");
  CFNumberRef vendorRef =
      IOHIDDeviceGetProperty(device, CFSTR(kIOHIDVendorIDKey));
  CFNumberRef productRef =
      IOHIDDeviceGetProperty(device, CFSTR(kIOHIDProductIDKey));
  int v = 0, p = 0;
  if (vendorRef)
    CFNumberGetValue(vendorRef, kCFNumberIntType, &v);
  if (productRef)
    CFNumberGetValue(productRef, kCFNumberIntType, &p);

  HIDDeviceContext *ctx = malloc(sizeof(HIDDeviceContext));
  ctx->device = device;
  ctx->vendorID = v;
  ctx->productID = p;
  ctx->manager = (IOHIDManagerRef)sender;
  ctx->runLoop = CFRunLoopGetCurrent();

  IOHIDDeviceRegisterInputValueCallback(device, HIDCallback, ctx);
  IOHIDDeviceScheduleWithRunLoop(device, CFRunLoopGetCurrent(),
                                 kCFRunLoopDefaultMode);
  IOHIDDeviceRegisterRemovalCallback(device, DeviceRemovalCallback, ctx);

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

bool setupDevice(IOHIDManagerRef *hidManager, int vendorID, int productID) {
  printf("setup device .... CFRunLoopGetCurrent() = %p\n",
         CFRunLoopGetCurrent());
  *hidManager = IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
  IOHIDManagerRegisterDeviceMatchingCallback(*hidManager,
                                             DeviceMatchingCallback, NULL);
  IOHIDManagerRegisterDeviceRemovalCallback(*hidManager,
                                            ManagerDeviceRemovalCallback, NULL);
  IOReturn ret = IOHIDManagerOpen(*hidManager, kIOHIDOptionsTypeNone);
  if (ret != kIOReturnSuccess) {
    printf("IOHIDManagerOpen failed: %d\n", ret);
    CFRelease(*hidManager);
    return false;
  }
  // setup listener device
  CFDictionaryRef dict = CreateMatchingDict(vendorID, productID);
  // use MatchingMultiple because SetDeviceMatching overrides prev matches
  IOHIDManagerSetDeviceMatchingMultiple(
      *hidManager,
      CFArrayCreate(NULL, (const void **)&dict, 1, &kCFTypeArrayCallBacks));
  // IOHIDManagerSetDeviceMatching(hidManager, dict);
  CFRelease(dict);
  // check  if device is available now
  CFSetRef deviceSet = IOHIDManagerCopyDevices(*hidManager);
  if (!deviceSet) {
    printf("No HID devices found.\n");
    CFRelease(*hidManager);
    return false;
  }

  if (!deviceSet || CFSetGetCount(deviceSet) == 0) {
    if (deviceSet) {
      CFRelease(deviceSet);
    }
    CFRelease(*hidManager);
    return false;
  }

  bool deviceFound = false;
  IOHIDDeviceRef *devices =
      malloc(sizeof(IOHIDDeviceRef) * CFSetGetCount(deviceSet));
  CFSetGetValues(deviceSet, (const void **)devices);
  for (CFIndex i = 0; i < CFSetGetCount(deviceSet); i++) {
    IOHIDDeviceRef dev = devices[i];

    CFNumberRef vendorRef =
        IOHIDDeviceGetProperty(dev, CFSTR(kIOHIDVendorIDKey));
    CFNumberRef productRef =
        IOHIDDeviceGetProperty(dev, CFSTR(kIOHIDProductIDKey));

    int vid = 0, pid = 0;
    if (vendorRef)
      CFNumberGetValue(vendorRef, kCFNumberIntType, &vid);
    if (productRef)
      CFNumberGetValue(productRef, kCFNumberIntType, &pid);

    if (vid == vendorID && pid == productID) {
      deviceFound = true;
    }
  }
  free(devices);
  CFRelease(deviceSet);

  if (deviceFound == false) {
    CFRelease(*hidManager);
    return false;
  }
  return true;
}

void Start(IOHIDManagerRef hidManager, CFRunLoopRef runLoop) {
  printf("Starting run loop...\n");
  printf("Start : = loop %p \n", runLoop);
  printf("hid %p\n", hidManager);
  IOHIDManagerScheduleWithRunLoop(hidManager, runLoop, kCFRunLoopDefaultMode);
  // CFRunLoopPerformBlock(*runLoop,
  // kCFRunLoopDefaultMode, ^{
  //   printf("Starting run loop...\n");
  //   CFRunLoopRun();
  // });
  // CFRunLoopWakeUp(*runLoop);
  CFRunLoopRun();
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

  // INFO: a keyboard may appear multiple times in deviceSet, why?
  //  Composite devices: many keyboards are composite HID devices, they expose
  //  multiple logical interfaces f.e. a keyboard, a consumer control (vol
  //  up/down, brightness, etc.), possibly a touchpad or mouse if integrated.

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
