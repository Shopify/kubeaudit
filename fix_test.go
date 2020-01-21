package kubeaudit_test

// const fixtureDir = "internal/test/fixtures"

// TODO Reenable once all auditors have been added
// func TestFix(t *testing.T) {
// 	cases := []struct {
// 		origFile  string
// 		fixedFile string
// 	}{
// 		{"all-resources_v1.yml", "all-resources-fixed_v1.yml"},
// 		{"all-resources_v1beta1.yml", "all-resources-fixed_v1beta1.yml"},
// 		{"all-resources_v1beta2.yml", "all-resources-fixed_v1beta2.yml"},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.origFile+" <=> "+tt.fixedFile, func(t *testing.T) {
// 			assert := assert.New(t)

// 			report := test.AuditManifest(t, fixtureDir, tt.origFile, all.Auditors())

// 			fixed := bytes.NewBuffer(nil)
// 			report.Fix(fixed)

// 			expected, err := ioutil.ReadFile(filepath.Join(fixtureDir, tt.fixedFile))
// 			assert.Nil(err)

// 			assert.Equal(string(expected), fixed.String())
// 		})
// 	}
// }
