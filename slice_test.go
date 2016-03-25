package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSliceContains_Success(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, "4")

	//verify
	require.Equal(t, true, r)
}

func TestSliceContains_False(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, "5")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_TypeMismatch(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, 4)

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_Nil(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, nil)

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_NotASlice(t *testing.T) {

	//setup

	//execute
	r := SliceContains("1", "1")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_Empty(t *testing.T) {

	//setup
	slice := []string{}

	//execute
	r := SliceContains(slice, "1")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_NilSlice(t *testing.T) {

	//setup

	//execute
	r := SliceContains(nil, "1")

	//verify
	require.Equal(t, false, r)
}

func TestStringMapToSlice_Success(t *testing.T) {

	//setup
	m := make(map[string]interface{})
	m["1"] = &Session{Id: "1", Version: 1}
	m["2"] = &Session{Id: "2", Version: 2}

	//execute
	r := StringMapToSlice(m)

	//verify
	sessions := r.([]Session)
	require.Equal(t, 2, len(sessions))
	require.Equal(t, true, sessions[0].Id == "1" || sessions[1].Id == "1")
	require.Equal(t, true, sessions[0].Id == "2" || sessions[1].Id == "2")
}

func TestChopSlice_Valid(t *testing.T) {

	//setup
	args := []Session{
		{Id: "1", UserId: "A"},
		{Id: "2", UserId: "B"},
		{Id: "3", UserId: "A"},
		{Id: "4", UserId: "B"},
		{Id: "5", UserId: "C"},
	}

	//execute
	r, err := ChopSlice(args, "UserId")

	//verify
	require.NoError(t, err)

	m := r.(map[string][]Session)
	require.Len(t, m["A"], 2)
	require.Equal(t, args[0], m["A"][0])
	require.Equal(t, args[2], m["A"][1])

	require.Len(t, m["B"], 2)
	require.Equal(t, args[1], m["B"][0])
	require.Equal(t, args[3], m["B"][1])

	require.Len(t, m["C"], 1)
	require.Equal(t, args[4], m["C"][0])
}

func TestChopSlice_ValidPointer(t *testing.T) {

	//setup
	args := []*Session{
		{Id: "1", UserId: "A"},
		{Id: "2", UserId: "B"},
		{Id: "3", UserId: "A"},
		{Id: "4", UserId: "B"},
		{Id: "5", UserId: "C"},
	}

	//execute
	r, err := ChopSlice(args, "UserId")

	//verify
	require.NoError(t, err)

	m := r.(map[string][]*Session)
	require.Len(t, m["A"], 2)
	require.Equal(t, args[0], m["A"][0])
	require.Equal(t, args[2], m["A"][1])

	require.Len(t, m["B"], 2)
	require.Equal(t, args[1], m["B"][0])
	require.Equal(t, args[3], m["B"][1])

	require.Len(t, m["C"], 1)
	require.Equal(t, args[4], m["C"][0])
}

func TestChopSlice_NilSlice(t *testing.T) {

	//setup

	//execute
	r, err := ChopSlice(nil, "UserId")

	//verify
	require.EqualError(t, err, "Slice is nil")
	require.Nil(t, r)
}

func TestChopSlice_EmptySlice(t *testing.T) {

	//setup
	args := []Session{}

	//execute
	r, err := ChopSlice(args, "UserId")

	//verify
	require.EqualError(t, err, "Slice is empty")
	require.Nil(t, r)
}

func TestChopSlice_NotSlice(t *testing.T) {

	//setup
	args := 0

	//execute
	r, err := ChopSlice(args, "UserId")

	//verify
	require.EqualError(t, err, "Argument is not a slice")
	require.Nil(t, r)
}

func TestChopSortedSlice_Valid(t *testing.T) {

	//setup
	args := []Session{
		{Id: "1", UserId: "A"},
		{Id: "3", UserId: "A"},
		{Id: "2", UserId: "B"},
		{Id: "4", UserId: "B"},
		{Id: "5", UserId: "C"},
	}

	//execute
	r, err := ChopSortedSlice(args, "UserId")

	//verify
	require.NoError(t, err)

	m := r.(map[string][]Session)
	require.Len(t, m["A"], 2)
	require.Equal(t, args[0], m["A"][0])
	require.Equal(t, args[1], m["A"][1])

	require.Len(t, m["B"], 2)
	require.Equal(t, args[2], m["B"][0])
	require.Equal(t, args[3], m["B"][1])

	require.Len(t, m["C"], 1)
	require.Equal(t, args[4], m["C"][0])
}

func TestChopSortedSlice_ValidPointer(t *testing.T) {

	//setup
	args := []*Session{
		{Id: "1", UserId: "A"},
		{Id: "3", UserId: "A"},
		{Id: "2", UserId: "B"},
		{Id: "4", UserId: "B"},
		{Id: "5", UserId: "C"},
	}

	//execute
	r, err := ChopSortedSlice(args, "UserId")

	//verify
	require.NoError(t, err)

	m := r.(map[string][]*Session)
	require.Len(t, m["A"], 2)
	require.Equal(t, args[0], m["A"][0])
	require.Equal(t, args[1], m["A"][1])

	require.Len(t, m["B"], 2)
	require.Equal(t, args[2], m["B"][0])
	require.Equal(t, args[3], m["B"][1])

	require.Len(t, m["C"], 1)
	require.Equal(t, args[4], m["C"][0])
}

func TestChopSortedSlice_NilSlice(t *testing.T) {

	//setup

	//execute
	r, err := ChopSortedSlice(nil, "UserId")

	//verify
	require.EqualError(t, err, "Slice is nil")
	require.Nil(t, r)
}

func TestChopSortedSlice_EmptySlice(t *testing.T) {

	//setup
	args := []Session{}

	//execute
	r, err := ChopSortedSlice(args, "UserId")

	//verify
	require.EqualError(t, err, "Slice is empty")
	require.Nil(t, r)
}

func TestChopSortedSlice_NotSlice(t *testing.T) {

	//setup
	args := 0

	//execute
	r, err := ChopSortedSlice(args, "UserId")

	//verify
	require.EqualError(t, err, "Argument is not a slice")
	require.Nil(t, r)
}

func BenchmarkChopSlice_Every10(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 10)
	b.ResetTimer()
	chopSliceBenchmark(testSet, b)
}

func BenchmarkChopSlice_Every100(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 100)
	b.ResetTimer()
	chopSliceBenchmark(testSet, b)
}

func BenchmarkChopSlice_Every1000(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 1000)
	b.ResetTimer()
	chopSliceBenchmark(testSet, b)
}

func BenchmarkChopSortedSlice_Every10(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 10)
	b.ResetTimer()
	chopSortedSliceBenchmark(testSet, b)
}

func BenchmarkChopSortedSlice_Every100(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 100)
	b.ResetTimer()
	chopSortedSliceBenchmark(testSet, b)
}

func BenchmarkChopSortedSlice_Every1000(b *testing.B) {
	testSet := buildBenchmarkSlice(10000, 1000)
	b.ResetTimer()
	chopSortedSliceBenchmark(testSet, b)
}

func buildBenchmarkSlice(num, every int) interface{} {
	var testSet []Session
	for i := 0; i < num; i++ {
		testSet = append(testSet, Session{Id: string(i), UserId: string(i / every)})
	}
	return testSet
}

func chopSliceBenchmark(testSet interface{}, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := ChopSlice(testSet, "UserId")
		if err != nil {
			panic(err)
		}
	}
}

func chopSortedSliceBenchmark(testSet interface{}, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := ChopSortedSlice(testSet, "UserId")
		if err != nil {
			panic(err)
		}
	}
}