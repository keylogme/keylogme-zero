#ifndef LISTENER_DARWIN_H
#define LISTENER_DARWIN_H

#include <IOKit/hid/IOHIDManager.h>
#include <stdbool.h>
#include <stdint.h>

bool setupDevice(IOHIDManagerRef *hid, int vendorID, int productID);
void stopDevice(CFRunLoopRef runLoop, IOHIDManagerRef manager);
void Start(IOHIDManagerRef hid, CFRunLoopRef runLoop);
extern void GoHandleKeyEvent(int code, int value, int vendorID, int productID);
extern void GoHandleDeviceEvent(int vendorID, int productID,
                                int connected); // 1=connected, 0=disconnected

void ListConnectedHIDDevices();
#endif
