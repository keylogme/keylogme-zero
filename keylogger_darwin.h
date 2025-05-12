#ifndef LISTENER_DARWIN_H
#define LISTENER_DARWIN_H

#include <stdint.h>

void setupDevice(int vendorID, int productID);
void Start(void);
void Stop(void);
extern void GoHandleKeyEvent(int code, int value, int vendorID, int productID);
extern void GoHandleDeviceEvent(int vendorID, int productID,
                                int connected); // 1=connected, 0=disconnected

void ListConnectedHIDDevices();
#endif
