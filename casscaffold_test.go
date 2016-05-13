package gowork

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestBuildUpdateStatement_Success(t *testing.T) {

	//setup
	cas := Cassandra{}

	config := Config{
		Id: "TEST",
		Created: CurrentTime(),
		Description: "descr",
		Value: "val",
		Version: 1,
	}

	overrides := make(map[string]interface{})

	//execute
	str, params := cas.BuildUpdateStatement("test", &config, overrides)

	//verify
	require.Equal(t, "UPDATE test SET cr = ?, descr = ?, value = ?, v = ? WHERE id = ?", str)
	require.Equal(t, 5, len(params))
	require.Equal(t, config.Created, params[0])
	require.Equal(t, "descr", params[1])
	require.Equal(t, "val", params[2])
	require.Equal(t, 1, params[3])
	require.Equal(t, "TEST", params[4])
}

func TestBuildUpdateStatement_Override(t *testing.T) {

	//setup
	cas := Cassandra{}

	config := Config{
		Id: "TEST",
		Created: CurrentTime(),
		Description: "descr",
		Value: "val",
		Version: 1,
	}

	overrides := map[string]interface{}{
		"descr": "BLAH",
	}

	//execute
	str, params := cas.BuildUpdateStatement("test", &config, overrides)

	//verify
	require.Equal(t, "UPDATE test SET cr = ?, descr = ?, value = ?, v = ? WHERE id = ?", str)
	require.Equal(t, 5, len(params))
	require.Equal(t, config.Created, params[0])
	require.Equal(t, "BLAH", params[1])
	require.Equal(t, "val", params[2])
	require.Equal(t, 1, params[3])
	require.Equal(t, "TEST", params[4])
}

func TestBuildUpdateStatement_WithPartition(t *testing.T) {

	//setup
	cas := Cassandra{
		Partition: "date",
	}

	config := Config{
		Id: "TEST",
		Created: CurrentTime(),
		Description: "descr",
		Value: "val",
		Version: 1,
	}

	overrides := map[string]interface{}{
		"date": FloorDay(CurrentTime()),
	}

	//execute
	str, params := cas.BuildUpdateStatement("test", &config, overrides)

	//verify
	require.Equal(t, "UPDATE test SET cr = ?, descr = ?, value = ?, v = ? WHERE id = ? AND date = ?", str)
	require.Equal(t, 6, len(params))
	require.Equal(t, config.Created, params[0])
	require.Equal(t, "descr", params[1])
	require.Equal(t, "val", params[2])
	require.Equal(t, 1, params[3])
	require.Equal(t, "TEST", params[4])
}
