package ets

import (
	"github.com/go-test/deep"
	"testing"
)

func TestVersion5_7_2_743(t *testing.T) {
	archive, err := OpenExportArchive("Testproject-5.7.2-743.knxproj", "")
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
					Line{ID: LineID("L-3"), ProjectID: ProjectID("P-0497-0"), Name: "New line", Address: 1, Devices: []DeviceInstance{}},
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

	if is, want := len(archive.ManufacturerFiles), 0; is != want {
		t.Fatalf("%v != %v", is, want)
	}
	if is, want := len(archive.HardwareFiles), 0; is != want {
		t.Fatalf("%v != %v", is, want)
	}
}
