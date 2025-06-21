#include "keylogger_darwin.h"
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/hid/IOHIDKeys.h>
#include <IOKit/hid/IOHIDManager.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static IOHIDManagerRef hidManager = NULL;
static CFRunLoopRef runLoop = NULL;
static bool runLoopStarted = false;

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
  if (!ctx || !element) {
    printf("HIDCallback: context or element is NULL\n");
    return;
  }
  uint32_t usagePage = IOHIDElementGetUsagePage(element);
  uint32_t usage = IOHIDElementGetUsage(element);
  CFIndex pressed = IOHIDValueGetIntegerValue(value);

  // printf("HIDCallback: usagePage: %d, usage: %d, pressed: %ld\n", usagePage,
  //        usage, pressed);
  if (usagePage == kHIDPage_KeyboardOrKeypad && usage >= 4 && usage <= 231) {
    // printf("Key usage: 0x%x, pressed: %ld\n", usage, pressed);
    GoHandleKeyEvent(usage, (int)pressed, ctx->vendorID, ctx->productID);
  } else if (usagePage == kHIDPage_Button && usage >= 1 && usage <= 8) {
    // printf("Click usage: 0x%x, pressed: %ld\n", usage, pressed);
    GoHandleKeyEvent(usage, (int)pressed, ctx->vendorID, ctx->productID);
  }
}

// void stopDevice(CFRunLoopRef runLoop, IOHIDManagerRef manager) {
//
//   IOHIDManagerUnscheduleFromRunLoop(manager, runLoop, kCFRunLoopDefaultMode);
//   IOHIDManagerClose(manager, kIOHIDOptionsTypeNone);
//   CFRelease(manager);
//
//   // Stop loop
//   CFRunLoopPerformBlock(runLoop, kCFRunLoopDefaultMode, ^{
//     CFRunLoopStop(runLoop);
//   });
//   CFRunLoopWakeUp(runLoop);
//   printf("Finished stopping device\n");
// }

void DeviceRemovalCallback(void *context, IOReturn result, void *sender) {
  HIDDeviceContext *ctx = (HIDDeviceContext *)context;
  if (!ctx)
    return;
  // printf("Device removal callback vendor %d | product %d\n", ctx->vendorID,
  //        ctx->productID);

  IOHIDDeviceUnscheduleFromRunLoop(ctx->device, ctx->runLoop,
                                   kCFRunLoopDefaultMode);
  IOHIDDeviceClose(ctx->device, kIOHIDOptionsTypeNone);

  GoHandleDeviceEvent(ctx->vendorID, ctx->productID, 0); // disconnected

  // cleanup
  // stopDevice(ctx->runLoop, ctx->manager);
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
  // printf("Device matching callback\n");
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

  IOReturn openResult = IOHIDDeviceOpen(device, kIOHIDOptionsTypeNone);
  if (openResult != kIOReturnSuccess) {
    printf("IOHIDDeviceOpen failed: %d\n", openResult);
    free(ctx);
    return;
  }
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

CFMutableDictionaryRef CreateMatchingKeyboard() {
  CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
      kCFAllocatorDefault, 0, &kCFTypeDictionaryKeyCallBacks,
      &kCFTypeDictionaryValueCallBacks);
  int usagePage = kHIDPage_GenericDesktop;
  CFNumberRef genericDesktop =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usagePage);
  int usage = kHIDUsage_GD_Keyboard;
  CFNumberRef keyboard =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usage);

  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsagePageKey), genericDesktop);
  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsageKey), keyboard);

  CFRelease(genericDesktop);
  CFRelease(keyboard);
  return dict;
}

CFMutableDictionaryRef CreateMatchingPointer() {
  CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
      kCFAllocatorDefault, 0, &kCFTypeDictionaryKeyCallBacks,
      &kCFTypeDictionaryValueCallBacks);
  int usagePage = kHIDPage_GenericDesktop;
  CFNumberRef genericDesktop =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usagePage);
  int usage = kHIDUsage_GD_Pointer;
  CFNumberRef pointer =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usage);

  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsagePageKey), genericDesktop);
  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsageKey), pointer);

  CFRelease(genericDesktop);
  CFRelease(pointer);
  return dict;
}

CFMutableDictionaryRef CreateMatchingMouse() {
  CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
      kCFAllocatorDefault, 0, &kCFTypeDictionaryKeyCallBacks,
      &kCFTypeDictionaryValueCallBacks);
  int usagePage = kHIDPage_GenericDesktop;
  CFNumberRef genericDesktop =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usagePage);
  int usage = kHIDUsage_GD_Mouse;
  CFNumberRef mouse =
      CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &usage);

  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsagePageKey), genericDesktop);
  CFDictionarySetValue(dict, CFSTR(kIOHIDDeviceUsageKey), mouse);

  CFRelease(genericDesktop);
  CFRelease(mouse);
  return dict;
}

CFMutableDictionaryRef
CreateMatchingDictByProductName(const char *productName) {
  if (productName == NULL) {
    return NULL;
  }

  // Convert C string to CFStringRef
  CFStringRef cfProductName = CFStringCreateWithCString(
      kCFAllocatorDefault, productName, kCFStringEncodingUTF8);
  if (!cfProductName) {
    return NULL;
  }

  // Create dictionary
  CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
      kCFAllocatorDefault, 1, &kCFTypeDictionaryKeyCallBacks,
      &kCFTypeDictionaryValueCallBacks);

  if (dict) {
    CFDictionarySetValue(dict, CFSTR(kIOHIDProductKey), cfProductName);
  }

  CFRelease(cfProductName);
  return dict;
}

void Start() {
  // printf("%d\n", runLoopStarted);
  if (hidManager) {
    // printf("Run loop already started, not setting up device again.\n");
    return;
  }
  // printf("setup device .... CFRunLoopGetCurrent() = %p\n",
  // CFRunLoopGetCurrent());
  hidManager = IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
  IOHIDManagerRegisterDeviceMatchingCallback(hidManager, DeviceMatchingCallback,
                                             NULL);
  IOHIDManagerRegisterDeviceRemovalCallback(hidManager,
                                            ManagerDeviceRemovalCallback, NULL);
  IOReturn ret = IOHIDManagerOpen(hidManager, kIOHIDOptionsTypeNone);
  if (ret != kIOReturnSuccess) {
    printf("IOHIDManagerOpen failed: %d\n", ret);
    CFRelease(hidManager);
    return;
  }
  // setup listener device
  CFMutableArrayRef matchingArray =
      CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
  CFDictionaryRef dictKeyboard = CreateMatchingKeyboard();
  CFArrayAppendValue(matchingArray, dictKeyboard);
  CFRelease(dictKeyboard);
  CFDictionaryRef dictPointer = CreateMatchingPointer();
  CFArrayAppendValue(matchingArray, dictPointer);
  CFRelease(dictPointer);
  CFDictionaryRef dictMouse = CreateMatchingMouse();
  CFArrayAppendValue(matchingArray, dictMouse);
  CFRelease(dictMouse);

  // INFO: use MatchingMultiple because SetDeviceMatching overrides prev
  // matches
  IOHIDManagerSetDeviceMatchingMultiple(hidManager, matchingArray);
  CFRelease(matchingArray);
  // Start
  runLoop = CFRunLoopGetCurrent();
  IOHIDManagerScheduleWithRunLoop(hidManager, runLoop, kCFRunLoopDefaultMode);
  runLoopStarted = true;

  // Blocking call
  // printf("Starting run loop...\n");
  CFRunLoopRun();

  // cleanup
  Stop();
  runLoopStarted = false;
  hidManager = NULL;
  runLoop = NULL;
}

void Stop() {
  if (!runLoopStarted) {
    return;
  }
  if (hidManager && runLoop) {
    IOHIDManagerUnscheduleFromRunLoop(hidManager, runLoop,
                                      kCFRunLoopDefaultMode);
    IOHIDManagerClose(hidManager, kIOHIDOptionsTypeNone);
    CFRelease(hidManager);
    hidManager = NULL;
    // printf("Run loop unscheduled and HID manager released.\n");
    CFRunLoopStop(runLoop);
    runLoop = NULL;
    runLoopStarted = false;
    // printf("Run loop stopped.\n");
  }
}

bool checkDeviceIsConnected(int vendorID, int productID) {
  IOHIDManagerRef mgr =
      IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
  IOHIDManagerOpen(mgr, kIOHIDOptionsTypeNone);

  // match
  CFMutableArrayRef matchingArray =
      CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
  CFDictionaryRef dictKeyboard = CreateMatchingKeyboard();
  CFArrayAppendValue(matchingArray, dictKeyboard);
  CFRelease(dictKeyboard);
  CFDictionaryRef dictPointer = CreateMatchingPointer();
  CFArrayAppendValue(matchingArray, dictPointer);
  CFRelease(dictPointer);
  CFDictionaryRef dictMouse = CreateMatchingMouse();
  CFArrayAppendValue(matchingArray, dictMouse);
  CFRelease(dictMouse);

  IOHIDManagerSetDeviceMatchingMultiple(mgr, matchingArray);
  CFRelease(matchingArray);
  //

  CFSetRef deviceSet = IOHIDManagerCopyDevices(mgr);
  if (!deviceSet) {
    printf("No HID devices found.\n");
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
      break;
    }
  }
  free(devices);
  CFRelease(deviceSet);
  CFRelease(mgr);

  if (deviceFound == false) {
    return false;
  }
  return true;
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
    printf("Device: %s | Product ID: 0x%04x | Vendor ID: 0x%04x | %d %d\n",
           name, pid, vid, pid, vid);
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
