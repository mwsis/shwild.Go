package shwild_test

import (
	shwild "github.com/synesissoftware/shwild.Go"

	"fmt"
	"path"
	"runtime"
	"testing"
)

/* /////////////////////////////////////////////////////////////////////////
 * internal functions
 */

func check_Match(t *testing.T, pattern, s string, expectedResult bool, e error) {

	m_r, m_e := shwild.Match(pattern, s)

	if expectedResult == m_r && e == m_e {

		return
	}

	_, file, line, hasCallInfo := runtime.Caller(1)

	var msg string

	if hasCallInfo {

		if expectedResult != m_r {

			msg = fmt.Sprintf("\t%s:%d: Match('%s', '%s') returned '%v'; '%v' expected", path.Base(file), line, pattern, s, m_r, expectedResult)
		}
	} else {

	}

	fmt.Printf("%s\n", msg)

	t.Fail()
}

/* /////////////////////////////////////////////////////////////////////////
 * tests
 */

func Test_Match_with_empty_pattern(t *testing.T) {

	check_Match(t, "", "", true, nil)
	check_Match(t, "", "1", false, nil)
	check_Match(t, "", "*", false, nil)
	check_Match(t, "", ".", false, nil)
}

func Test_Match_with_wild1(t *testing.T) {

	check_Match(t, "?", "", false, nil)
	check_Match(t, "?", "?", true, nil)
}

func Test_Match_with_allstar_patterns(t *testing.T) {

	// 1 star

	check_Match(t, "*", "", true, nil)
	check_Match(t, "*", "1", true, nil)
	check_Match(t, "*", "*", true, nil)
	check_Match(t, "*", ".", true, nil)

	// 2 star

	check_Match(t, "**", "", true, nil)
	check_Match(t, "**", "1", true, nil)
	check_Match(t, "**", "*", true, nil)
	check_Match(t, "**", ".", true, nil)

	// 5 star

	check_Match(t, "*****", "", true, nil)
	check_Match(t, "*****", "1", true, nil)
	check_Match(t, "*****", "*", true, nil)
	check_Match(t, "*****", ".", true, nil)
}

func Test_Match_with_literal(t *testing.T) {

	check_Match(t, "a", "a", true, nil)
	check_Match(t, "aa", "a", false, nil)
	check_Match(t, "aa", "aa", true, nil)
	check_Match(t, "a", "aa", false, nil)
}

func Test_Match_with_literal_and_wild1(t *testing.T) {

	check_Match(t, "a?", "a", false, nil)
	check_Match(t, "a?", "a?", true, nil)
	check_Match(t, "a?", "aa", true, nil)
	check_Match(t, "a?", "aaa", false, nil)
	check_Match(t, "?a", "a", false, nil)
	check_Match(t, "?a", "a?", false, nil)
	check_Match(t, "?a", "?a", true, nil)
	check_Match(t, "?a", "aa", true, nil)
	check_Match(t, "?a", "aaa", false, nil)
}

func Test_Match_with_literal_and_wildN(t *testing.T) {

	check_Match(t, "a*", "a", true, nil)
	check_Match(t, "a*", "a*", true, nil)
	check_Match(t, "a*", "aa", true, nil)
	check_Match(t, "a*", "aaa", true, nil)
	check_Match(t, "a*", "abcdefghijklmno", true, nil)
	check_Match(t, "a*o", "abcdefghijklmno", true, nil)
	check_Match(t, "a*n", "abcdefghijklmno", false, nil)
	check_Match(t, "a*n*", "abcdefghijklmno", true, nil)
	check_Match(t, "*a", "a", true, nil)
	check_Match(t, "*a", "a*", false, nil)
	check_Match(t, "*a", "*a", true, nil)
	check_Match(t, "*a", "aa", true, nil)
	check_Match(t, "*a", "aaa", true, nil)
}

func Test_Match_with_explicit_range(t *testing.T) {

	check_Match(t, "[abc]", "a", true, nil)
	check_Match(t, "[abc]", "b", true, nil)
	check_Match(t, "[abc]", "c", true, nil)
	check_Match(t, "[abc]", "d", false, nil)

	check_Match(t, "[-abc]", "-", true, nil)
}

func Test_Match_with_forward_continuum_range(t *testing.T) {

	check_Match(t, "[a-c]", "a", true, nil)
	check_Match(t, "[a-c]", "b", true, nil)
	check_Match(t, "[a-c]", "c", true, nil)
	check_Match(t, "[a-c]", "-", false, nil)
	check_Match(t, "[a-c]", "z", false, nil)
	check_Match(t, "[a-c]", "A", false, nil)
	check_Match(t, "[a-c]", "B", false, nil)
	check_Match(t, "[a-c]", "C", false, nil)
	check_Match(t, "[a-c]", "D", false, nil)
	check_Match(t, "[a-c]", "E", false, nil)

	check_Match(t, "[-ac]", "a", true, nil)
	check_Match(t, "[-ac]", "b", false, nil)
	check_Match(t, "[-ac]", "c", true, nil)
	check_Match(t, "[-ac]", "-", true, nil)
	check_Match(t, "[a-c]", "d", false, nil)
	check_Match(t, "[a-c]", "z", false, nil)
}

func Test_Match_with_forward_continuum_notrange(t *testing.T) {

	check_Match(t, "[^a-c]", "a", false, nil)
	check_Match(t, "[^a-c]", "b", false, nil)
	check_Match(t, "[^a-c]", "c", false, nil)
	check_Match(t, "[^a-c]", "-", true, nil)
	check_Match(t, "[^a-c]", "z", true, nil)
	check_Match(t, "[^a-c]", "A", true, nil)
	check_Match(t, "[^a-c]", "B", true, nil)
	check_Match(t, "[^a-c]", "C", true, nil)
	check_Match(t, "[^a-c]", "D", true, nil)
	check_Match(t, "[^a-c]", "E", true, nil)

	check_Match(t, "[^-ac]", "a", false, nil)
	check_Match(t, "[^-ac]", "b", true, nil)
	check_Match(t, "[^-ac]", "c", false, nil)
	check_Match(t, "[^-ac]", "-", false, nil)
	check_Match(t, "[^a-c]", "d", true, nil)
	check_Match(t, "[^a-c]", "z", true, nil)
}

func Test_Match_with_backward_continuum_range(t *testing.T) {

	check_Match(t, "[c-a]", "a", true, nil)
	check_Match(t, "[c-a]", "b", true, nil)
	check_Match(t, "[c-a]", "c", true, nil)
	check_Match(t, "[c-a]", "-", false, nil)
	check_Match(t, "[c-a]", "z", false, nil)
	check_Match(t, "[c-a]", "A", false, nil)
	check_Match(t, "[c-a]", "B", false, nil)
	check_Match(t, "[c-a]", "C", false, nil)
	check_Match(t, "[c-a]", "D", false, nil)
	check_Match(t, "[c-a]", "E", false, nil)
}

func Test_Match_with_backward_continuum_notrange(t *testing.T) {

	check_Match(t, "[^c-a]", "a", false, nil)
	check_Match(t, "[^c-a]", "b", false, nil)
	check_Match(t, "[^c-a]", "c", false, nil)
	check_Match(t, "[^c-a]", "-", true, nil)
	check_Match(t, "[^c-a]", "z", true, nil)
	check_Match(t, "[^c-a]", "A", true, nil)
	check_Match(t, "[^c-a]", "B", true, nil)
	check_Match(t, "[^c-a]", "C", true, nil)
	check_Match(t, "[^c-a]", "D", true, nil)
	check_Match(t, "[^c-a]", "E", true, nil)
}

func Test_Match_with_forward_crosscase_continuum_range(t *testing.T) {

	check_Match(t, "[a-C]", "-", false, nil)
	check_Match(t, "[a-C]", "a", true, nil)
	check_Match(t, "[a-C]", "b", true, nil)
	check_Match(t, "[a-C]", "c", true, nil)
	check_Match(t, "[a-C]", "d", false, nil)
	check_Match(t, "[a-C]", "z", false, nil)
	check_Match(t, "[a-C]", "A", true, nil)
	check_Match(t, "[a-C]", "B", true, nil)
	check_Match(t, "[a-C]", "C", true, nil)
	check_Match(t, "[a-C]", "D", false, nil)
	check_Match(t, "[a-C]", "E", false, nil)
}

func Test_Match_with_forward_crosscase_continuum_notrange(t *testing.T) {

	check_Match(t, "[^a-C]", "-", true, nil)
	check_Match(t, "[^a-C]", "a", false, nil)
	check_Match(t, "[^a-C]", "b", false, nil)
	check_Match(t, "[^a-C]", "c", false, nil)
	check_Match(t, "[^a-C]", "d", true, nil)
	check_Match(t, "[^a-C]", "z", true, nil)
	check_Match(t, "[^a-C]", "A", false, nil)
	check_Match(t, "[^a-C]", "B", false, nil)
	check_Match(t, "[^a-C]", "C", false, nil)
	check_Match(t, "[^a-C]", "D", true, nil)
	check_Match(t, "[^a-C]", "E", true, nil)
}

func Test_Match_with_escaped_special_characters(t *testing.T) {

	check_Match(t, "a\\*c", "a_c", false, nil)
	check_Match(t, "a\\*c", "a*c", true, nil)

	check_Match(t, "a\\?c", "a_c", false, nil)
	check_Match(t, "a\\?c", "a?c", true, nil)

	check_Match(t, "a\\[c", "a_c", false, nil)
	check_Match(t, "a\\[c", "a[c", true, nil)

	check_Match(t, "a\\]c", "a_c", false, nil)
	check_Match(t, "a\\]c", "a]c", true, nil)
	check_Match(t, "a]c", "a]c", true, nil)
}

func Test_Pattern_Match_from_examples_1(t *testing.T) {

	pattern := "[ER]*"

	check_Match(t, pattern, "", false, nil)
	check_Match(t, pattern, "E", true, nil)
	check_Match(t, pattern, "EX", true, nil)
	check_Match(t, pattern, "EXAMPLES.md", true, nil)
	check_Match(t, pattern, "README.md", true, nil)
	check_Match(t, pattern, "LICENSE", false, nil)
}

/* ///////////////////////////// end of file //////////////////////////// */
