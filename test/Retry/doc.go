// Retry if a source Observable sends an error notification, resubscribe to it
// in the hopes that it will complete without error. If count is zero or
// negative, the retry count will be effectively infinite. The scheduler
// passed when subscribing is used by Retry to schedule any retry attempt. The
// time between retries is 1 millisecond, so retry frequency is 1 kHz. Any
// SubscribeOn operators should be called after Retry to prevent lockups
// caused by mixing different schedulers in the same subscription for retrying
// and subscribing.
package Retry
