#ifndef LISTENER_DARWIN_H
#define LISTENER_DARWIN_H

#include <IOKit/hid/IOHIDManager.h>
#include <stdbool.h>
#include <stdint.h>

bool checkDeviceIsConnected(int vendorID, int productID);
void Stop();
void Start();
extern void GoHandleKeyEvent(int code, int value, int vendorID, int productID);
extern void GoHandleDeviceEvent(int vendorID, int productID,
                                int connected); // 1=connected, 0=disconnected

void ListConnectedHIDDevices();
#endif
