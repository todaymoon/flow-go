// Code generated by "stringer -type=initializerKind"; DO NOT EDIT.

package sema

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[initializerKindUnknown-0]
	_ = x[initializerKindComposite-1]
	_ = x[initializerKindInterface-2]
}

const _initializerKind_name = "initializerKindUnknowninitializerKindCompositeinitializerKindInterface"

var _initializerKind_index = [...]uint8{0, 22, 46, 70}

func (i initializerKind) String() string {
	if i < 0 || i >= initializerKind(len(_initializerKind_index)-1) {
		return "initializerKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _initializerKind_name[_initializerKind_index[i]:_initializerKind_index[i+1]]
}
