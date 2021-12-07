package ets

import (
	"github.com/go-test/deep"
	"testing"
)

func TestVersion6_0(t *testing.T) {
	archive, err := OpenExportArchive("Testproject-6.0.knxproj", "")
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()

	if is, want := len(archive.ProjectFiles), 1; is != want {
		t.Fatalf("%v != %v", is, want)
	}

	fproj := archive.ProjectFiles[0]
	t.Run("project", func(t *testing.T) {
		projInfo, err := fproj.Decode()
		if err != nil {
			t.Fatal(err)
		}

		if is, want := projInfo.Name, "Testproject"; is != want {
			t.Fatalf("%v != %v", is, want)
		}
		if is, want := projInfo.Comment, ""; is != want {
			t.Fatalf("%v != %v", is, want)
		}

		if is, want := len(fproj.InstallationFiles), 1; is != want {
			t.Fatalf("%v != %v", is, want)
		}

		finst := fproj.InstallationFiles[0]
		proj, err := finst.Decode()
		if err != nil {
			t.Fatal(err)
		}

		if is, want := len(proj.Installations), 1; is != want {
			t.Fatalf("%v != %v", is, want)
		}

		inst := proj.Installations[0]
		areas := []Area{
			Area{
				ID:        AreaID("A-1"),
				ProjectID: ProjectID("P-0497-0"),
				Name:      "Backbone area",
				Address:   0,
				Lines: []Line{
					Line{ID: LineID("L-1"), ProjectID: ProjectID("P-0497-0"), Name: "Backbone line", Address: 0, Devices: []DeviceInstance{}},
				},
			},
			Area{
				ID:        AreaID("A-2"),
				ProjectID: ProjectID("P-0497-0"),
				Name:      "New area",
				Address:   1,
				Lines: []Line{
					Line{ID: LineID("L-2"), ProjectID: ProjectID("P-0497-0"), Name: "Main line", Address: 0, Devices: []DeviceInstance{}},
					Line{ID: LineID("L-3"), ProjectID: ProjectID("P-0497-0"), Name: "New line", Address: 1, Devices: []DeviceInstance{
						DeviceInstance{
							ID:                 DeviceInstanceID("DI-1"),
							ProjectID:          ProjectID("P-0497-0"),
							ManufacturerID:     ManufacturerID("M-0007"),
							HardwareID:         HardwareID("H-6131.2F20-1"),
							ProductID:          ProductID("P-6131.2F20"),
							Hardware2ProgramID: Hardware2ProgramID("HP-3120-32-269B-3120-42-4C77"),
							Name:               "",
							Address:            1,
							ComObjects: []ComObjectInstanceRef{
								ComObjectInstanceRef{
									ComObjectRefID: ComObjectRefID("R-1"),
									ComObjectID:    ComObjectID("O-10"),
									DatapointType:  "",
									Links:          []string{"GA-1"},
								},
							},
						},
						DeviceInstance{
							ID:                 DeviceInstanceID("DI-2"),
							ProjectID:          ProjectID("P-0497-0"),
							ManufacturerID:     ManufacturerID("M-0083"),
							HardwareID:         HardwareID("H-4-2"),
							ProductID:          ProductID("P-AMS.2D1216.2E02"),
							Hardware2ProgramID: Hardware2ProgramID("HP-0019-21-D29E"),
							Name:               "",
							Address:            2,
							ComObjects: []ComObjectInstanceRef{
								ComObjectInstanceRef{
									ComObjectRefID: ComObjectRefID("R-10000"),
									ComObjectID:    ComObjectID("O-0"),
									DatapointType:  "",
									Links:          []string{"GA-1"},
								},
							},
						},
					}},
				},
			},
		}

		if diff := deep.Equal(inst.Topology, areas); diff != nil {
			t.Error(diff)
		}

		locations := []Space{
			Space{
				ID:                SpaceID("BP-1"),
				ProjectID:         ProjectID("P-0497-0"),
				DeviceInstanceIDs: []DeviceInstanceID{},
				Type:              "Building",
				Name:              "Testproject",
				SubSpaces: []Space{
					Space{
						ID:                SpaceID("BP-2"),
						ProjectID:         ProjectID("P-0497-0"),
						DeviceInstanceIDs: []DeviceInstanceID{},
						Type:              "BuildingPart",
						Name:              "Indoor",
						SubSpaces: []Space{
							Space{
								ID:                SpaceID("BP-7"),
								ProjectID:         ProjectID("P-0497-0"),
								DeviceInstanceIDs: []DeviceInstanceID{},
								Type:              "Floor",
								Name:              "Basement",
								SubSpaces: []Space{
									Space{
										ID:                SpaceID("BP-4"),
										ProjectID:         ProjectID("P-0497-0"),
										DeviceInstanceIDs: []DeviceInstanceID{},
										Type:              "Room",
										Name:              "Kitchen",
										SubSpaces:         []Space{},
									},
									Space{
										ID:                SpaceID("BP-9"),
										ProjectID:         ProjectID("P-0497-0"),
										DeviceInstanceIDs: []DeviceInstanceID{},
										Type:              "Room",
										Name:              "Living room",
										SubSpaces:         []Space{},
									},
									Space{
										ID:                SpaceID("BP-11"),
										ProjectID:         ProjectID("P-0497-0"),
										DeviceInstanceIDs: []DeviceInstanceID{},
										Type:              "Corridor",
										Name:              "Corridor",
										SubSpaces: []Space{
											Space{
												ID:                SpaceID("BP-13"),
												ProjectID:         ProjectID("P-0497-0"),
												DeviceInstanceIDs: []DeviceInstanceID{},
												Type:              "DistributionBoard",
												Name:              "Distribution Board",
												SubSpaces:         []Space{},
											},
										},
									},
								},
							},
							Space{
								ID:                SpaceID("BP-8"),
								ProjectID:         ProjectID("P-0497-0"),
								DeviceInstanceIDs: []DeviceInstanceID{},
								Type:              "Floor",
								Name:              "Upstairs",
								SubSpaces: []Space{
									Space{
										ID:                SpaceID("BP-5"),
										ProjectID:         ProjectID("P-0497-0"),
										DeviceInstanceIDs: []DeviceInstanceID{},
										Type:              "Room",
										Name:              "Bedroom",
										SubSpaces:         []Space{},
									},
									Space{
										ID:                SpaceID("BP-10"),
										ProjectID:         ProjectID("P-0497-0"),
										DeviceInstanceIDs: []DeviceInstanceID{},
										Type:              "Room",
										Name:              "Bathroom",
										SubSpaces:         []Space{},
									},
								},
							},
							Space{
								ID:                SpaceID("BP-12"),
								ProjectID:         ProjectID("P-0497-0"),
								DeviceInstanceIDs: []DeviceInstanceID{},
								Type:              "Stairway",
								Name:              "Stairway",
								SubSpaces:         []Space{},
							},
						},
					},
					Space{
						ID:                SpaceID("BP-3"),
						ProjectID:         ProjectID("P-0497-0"),
						DeviceInstanceIDs: []DeviceInstanceID{},
						Type:              "BuildingPart",
						Name:              "Outdoor",
						SubSpaces: []Space{
							Space{
								ID:                SpaceID("BP-6"),
								ProjectID:         ProjectID("P-0497-0"),
								DeviceInstanceIDs: []DeviceInstanceID{},
								Type:              "Room",
								Name:              "Garden",
								SubSpaces:         []Space{},
							},
						},
					},
				},
			},
		}
		if diff := deep.Equal(inst.Locations, locations); diff != nil {
			t.Error(diff)
		}
	})

	if is, want := len(archive.ManufacturerFiles), 3; is != want {
		t.Fatalf("%v != %v", is, want)
	}
	if is, want := len(archive.HardwareFiles), 2; is != want {
		t.Fatalf("%v != %v", is, want)
	}
}
