package rx

// Revision1234 is the revision identifier that changes whenever the rx library gets updated for every (even minor) change.
// This is because the library is used by copying fragments of the source of the library through jig instead of using a
// compiled version. When compiling the program that uses an older revision of rx, the build will fail and the user will
// know that it needs to run jig --regen to regenerate.
func Revision1234() {}
