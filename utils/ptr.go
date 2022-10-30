package utils

import "time"

func IntPtr(v int) *int                          { return &v }
func Int16Ptr(v int16) *int16                    { return &v }
func Int32Ptr(v int32) *int32                    { return &v }
func Int64Ptr(v int64) *int64                    { return &v }
func Float64Ptr(v float64) *float64              { return &v }
func BoolPtr(v bool) *bool                       { return &v }
func StringPtr(v string) *string                 { return &v }
func DurationPtr(v time.Duration) *time.Duration { return &v }

func IntPtrSafe(v *int, defaultValue int) *int {
	if v == nil {
		return &defaultValue
	}

	return v
}

func Int16PtrSafe(v *int16, defaultValue int16) *int16 {
	if v == nil {
		return &defaultValue
	}

	return v
}

func Int32PtrSafe(v *int32, defaultValue int32) *int32 {
	if v == nil {
		return &defaultValue
	}

	return v
}

func Int64PtrSafe(v *int64, defaultValue int64) *int64 {
	if v == nil {
		return &defaultValue
	}

	return v
}

func Float64PtrSafe(v *float64, defaultValue float64) *float64 {
	if v == nil {
		return &defaultValue
	}

	return v
}

func BoolPtrSafe(v *bool, defaultValue bool) *bool {
	if v == nil {
		return &defaultValue
	}

	return v
}

func StringPtrSafe(v *string, defaultValue string) *string {
	if v == nil {
		return &defaultValue
	}

	return v
}

func DurationPtrSafe(v *time.Duration, defaultValue time.Duration) *time.Duration {
	if v == nil {
		return &defaultValue
	}

	return v
}
