package typemap

import (
	"github.com/pkg/errors"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func FileComments(file *descriptor.FileDescriptorProto) (DefinitionComments, error) {
	return commentsAtPath([]int32{packagePath}, file), nil
}

func ServiceComments(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto) (DefinitionComments, error) {
	for i, s := range file.Service {
		if s == svc {
			path := []int32{servicePath, int32(i)}
			return commentsAtPath(path, file), nil
		}
	}
	return DefinitionComments{}, errors.Errorf("service not found in file")
}

func MethodComments(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto, method *descriptor.MethodDescriptorProto) (DefinitionComments, error) {
	for i, s := range file.Service {
		if s == svc {
			path := []int32{servicePath, int32(i)}
			for j, m := range s.Method {
				if m == method {
					path = append(path, serviceMethodPath, int32(j))
					return commentsAtPath(path, file), nil
				}
			}
		}
	}
	return DefinitionComments{}, errors.Errorf("service not found in file")
}

// DefinitionComments contains the comments surrounding a definition in a
// protobuf file.
//
// These follow the rules described by protobuf:
//
// A series of line comments appearing on consecutive lines, with no other
// tokens appearing on those lines, will be treated as a single comment.
//
// leading_detached_comments will keep paragraphs of comments that appear
// before (but not connected to) the current element. Each paragraph,
// separated by empty lines, will be one comment element in the repeated
// field.
//
// Only the comment content is provided; comment markers (e.g. //) are
// stripped out.  For block comments, leading whitespace and an asterisk
// will be stripped from the beginning of each line other than the first.
// Newlines are included in the output.
//
// Examples:
//
//   optional int32 foo = 1;  // Comment attached to foo.
//   // Comment attached to bar.
//   optional int32 bar = 2;
//
//   optional string baz = 3;
//   // Comment attached to baz.
//   // Another line attached to baz.
//
//   // Comment attached to qux.
//   //
//   // Another line attached to qux.
//   optional double qux = 4;
//
//   // Detached comment for corge. This is not leading or trailing comments
//   // to qux or corge because there are blank lines separating it from
//   // both.
//
//   // Detached comment for corge paragraph 2.
//
//   optional string corge = 5;
//   /* Block comment attached
//    * to corge.  Leading asterisks
//    * will be removed. */
//   /* Block comment attached to
//    * grault. */
//   optional int32 grault = 6;
//
//   // ignored detached comments.
type DefinitionComments struct {
	Leading         string
	Trailing        string
	LeadingDetached []string
}

func commentsAtPath(path []int32, sourceFile *descriptor.FileDescriptorProto) DefinitionComments {
	if sourceFile.SourceCodeInfo == nil {
		// The compiler didn't provide us with comments.
		return DefinitionComments{}
	}

	for _, loc := range sourceFile.SourceCodeInfo.Location {
		if pathEqual(path, loc.Path) {
			return DefinitionComments{
				Leading:         loc.GetLeadingComments(),
				LeadingDetached: loc.GetLeadingDetachedComments(),
				Trailing:        loc.GetTrailingComments(),
			}
		}
	}
	return DefinitionComments{}
}

func pathEqual(path1, path2 []int32) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i, v := range path1 {
		if path2[i] != v {
			return false
		}
	}
	return true
}