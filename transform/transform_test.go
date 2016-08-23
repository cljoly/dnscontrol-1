package transform

import (
	"net"
	"testing"
)

func TestIPToUint(t *testing.T) {
	ip := net.ParseIP("1.2.3.4")
	u, err := ipToUint(ip)
	if err != nil {
		t.Fatal(err)
	}
	if u != 16909060 {
		t.Fatalf("I to uint conversion failed. Should be 16909060. Got %d", u)
	}
	ip2 := UintToIP(u)
	if !ip.Equal(ip2) {
		t.Fatalf("IPs should be equal. %s is not %s", ip2, ip)
	}
}

func Test_DecodeTransformTable_failures(t *testing.T) {
	result, err := DecodeTransformTable("1.2.3.4 ~ 3.4.5.6")
	if result != nil {
		t.Errorf("expected nil, got (%v)\n", result)
	}
	if err == nil {
		t.Error("expect error, got none")
	}
}

func test_ip(t *testing.T, test string, expected string, actual net.IP) {
	if !net.ParseIP(expected).Equal(actual) {
		t.Errorf("Test %v: expected Low (%v), got (%v)\n", test, actual, expected)
	}
}

func Test_DecodeTransformTable_0(t *testing.T) {
	result, err := DecodeTransformTable("1.2.3.4 ~ 2.3.4.5 ~ 3.4.5.6 ~ ")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Test %v: expected col length (%v), got (%v)\n", 1, 1, len(result))
	}
	test_ip(t, "low", "1.2.3.4", result[0].Low)
	test_ip(t, "high", "2.3.4.5", result[0].High)
	test_ip(t, "newBase", "3.4.5.6", result[0].NewBase)
	//test_ip(t, "newIP", "", result[0].NewIPs)
}

func Test_DecodeTransformTable_1(t *testing.T) {
	result, err := DecodeTransformTable("1.2.3.4~2.3.4.5~3.4.5.6 ~;8.7.6.5 ~ 9.8.7.6 ~ 7.6.5.4 ~ ")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("Test %v: expected col length (%v), got (%v)\n", 1, 2, len(result))
	}
	test_ip(t, "Low[0]", "1.2.3.4", result[0].Low)
	test_ip(t, "High[0]", "2.3.4.5", result[0].High)
	test_ip(t, "NewBase[0]", "3.4.5.6", result[0].NewBase)
	//test_ip(t, "newIP[0]", "", result[0].NewIP)
	test_ip(t, "Low[1]", "8.7.6.5", result[1].Low)
	test_ip(t, "High[1]", "9.8.7.6", result[1].High)
	test_ip(t, "NewBase[1]", "7.6.5.4", result[1].NewBase)
	//test_ip(t, "newIP[1]", "", result[0].NewIP)
}
func Test_DecodeTransformTable_NewIP(t *testing.T) {
	result, err := DecodeTransformTable("1.2.3.4 ~ 2.3.4.5 ~  ~ 3.4.5.6 ")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Test %v: expected col length (%v), got (%v)\n", 1, 1, len(result))
	}
	test_ip(t, "low", "1.2.3.4", result[0].Low)
	test_ip(t, "high", "2.3.4.5", result[0].High)
	//test_ip(t, "newIP", "3.4.5.6", result[0].NewIP)
	test_ip(t, "newBase", "", result[0].NewBase)
}

func Test_DecodeTransformTable_order(t *testing.T) {
	raw := "9.8.7.6 ~ 8.7.6.5 ~ 7.6.5.4 ~"
	result, err := DecodeTransformTable(raw)
	if result != nil {
		t.Errorf("Invalid range not detected: (%v)\n", raw)
	}
	if err == nil {
		t.Error("expect error, got none")
	}
}

func Test_DecodeTransformTable_Base_and_IP(t *testing.T) {
	raw := "1.1.1.1~ 8.7.6.5 ~ 7.6.5.4 ~ 4.4.4.4"
	result, err := DecodeTransformTable(raw)
	if result != nil {
		t.Errorf("NewBase and NewIP should not both be specified: (%v)\n", raw)
	}
	if err == nil {
		t.Error("expect error, got none")
	}
}

func Test_TransformIP(t *testing.T) {

	var transforms1 = []IpConversion{{
		Low:     net.ParseIP("11.11.11.0"),
		High:    net.ParseIP("11.11.11.20"),
		NewBase: net.ParseIP("99.99.99.0"),
	}, {
		Low:     net.ParseIP("22.22.22.0"),
		High:    net.ParseIP("22.22.22.40"),
		NewBase: net.ParseIP("99.99.99.100"),
	}, {
		Low:     net.ParseIP("33.33.33.20"),
		High:    net.ParseIP("33.33.35.40"),
		NewBase: net.ParseIP("100.100.100.0"),
	}, {
		Low:     net.ParseIP("44.44.44.20"),
		High:    net.ParseIP("44.44.44.40"),
		NewBase: net.ParseIP("100.100.100.40"),
	}}

	var tests = []struct {
		experiment string
		expected   string
	}{
		{"11.11.11.0", "99.99.99.0"},
		{"11.11.11.1", "99.99.99.1"},
		{"11.11.11.11", "99.99.99.11"},
		{"11.11.11.19", "99.99.99.19"},
		{"11.11.11.20", "99.99.99.20"},
		{"11.11.11.21", "11.11.11.21"},
		{"22.22.22.22", "99.99.99.122"},
		{"22.22.22.255", "22.22.22.255"},
		{"33.33.33.0", "33.33.33.0"},
		{"33.33.33.19", "33.33.33.19"},
		{"33.33.33.20", "100.100.100.0"},
		{"33.33.33.21", "100.100.100.1"},
		{"33.33.33.33", "100.100.100.13"},
		{"33.33.35.39", "100.100.102.19"},
		{"33.33.35.40", "100.100.102.20"},
		{"33.33.35.41", "33.33.35.41"},
		{"44.44.44.24", "100.100.100.44"},
		{"44.44.44.44", "44.44.44.44"},
	}

	for _, test := range tests {
		experiment := net.ParseIP(test.experiment)
		expected := net.ParseIP(test.expected)
		actual, err := TransformIP(experiment, transforms1)
		if err != nil {
			t.Errorf("%v: got an err: %v\n", experiment, err)
		}
		if !expected.Equal(actual) {
			t.Errorf("%v: expected (%v) got (%v)\n", experiment, expected, actual)
		}
	}
}

func Test_TransformIP_NewIP(t *testing.T) {

	var transforms1 = []IpConversion{{
		Low:    net.ParseIP("11.11.11.0"),
		High:   net.ParseIP("11.11.11.20"),
		NewIPs: []net.IP{net.ParseIP("1.1.1.1")},
	}, {
		Low:    net.ParseIP("22.22.22.0"),
		High:   net.ParseIP("22.22.22.40"),
		NewIPs: []net.IP{net.ParseIP("2.2.2.2")},
	}, {
		Low:    net.ParseIP("33.33.33.20"),
		High:   net.ParseIP("33.33.35.40"),
		NewIPs: []net.IP{net.ParseIP("3.3.3.3")},
	},
	}

	var tests = []struct {
		experiment string
		expected   string
	}{
		{"11.11.11.0", "1.1.1.1"},
		{"11.11.11.1", "1.1.1.1"},
		{"11.11.11.11", "1.1.1.1"},
		{"11.11.11.19", "1.1.1.1"},
		{"11.11.11.20", "1.1.1.1"},
		{"11.11.11.21", "11.11.11.21"},
		{"22.22.22.22", "2.2.2.2"},
		{"22.22.22.255", "22.22.22.255"},
		{"33.33.33.0", "33.33.33.0"},
		{"33.33.33.19", "33.33.33.19"},
		{"33.33.33.20", "3.3.3.3"},
		{"33.33.33.21", "3.3.3.3"},
		{"33.33.33.33", "3.3.3.3"},
		{"33.33.35.39", "3.3.3.3"},
		{"33.33.35.40", "3.3.3.3"},
		{"33.33.35.41", "33.33.35.41"},
	}

	for _, test := range tests {
		experiment := net.ParseIP(test.experiment)
		expected := net.ParseIP(test.expected)
		actual, err := TransformIP(experiment, transforms1)
		if err != nil {
			t.Errorf("%v: got an err: %v\n", experiment, err)
		}
		if !expected.Equal(actual) {
			t.Errorf("%v: expected (%v) got (%v)\n", experiment, expected, actual)
		}
	}
}
