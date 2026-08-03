#import <ServiceManagement/ServiceManagement.h>
#import <Foundation/Foundation.h>

void registerLoginItemObjC(void) {
    if (@available(macOS 13.0, *)) {
        SMAppService *service = [SMAppService mainAppService];
        NSError *error = nil;
        BOOL ok = [service registerAndReturnError:&error];
        if (!ok && error != nil) {
            NSLog(@"[espresso] login item registration failed: %@", error.localizedDescription);
        }
    }
}
