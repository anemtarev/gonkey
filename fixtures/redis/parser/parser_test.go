package parser

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRedisFixtureParser_Load(t *testing.T) {
	type args struct {
		fixtures []string
	}

	type want struct {
		fixtures []*Fixture
		ctx      *context
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test basic",
			args: args{
				fixtures: []string{"redis"},
			},
			want: want{
				fixtures: []*Fixture{
					{
						Databases: map[int]Database{
							1: {
								Keys: &Keys{
									Values: map[string]*KeyValue{
										"key1": {
											Value: "value1",
										},
										"key2": {
											Value:      "value2",
											Expiration: time.Second * 10,
										},
									},
								},
								Sets: &Sets{
									Values: map[string]*SetRecordValue{
										"set1": {
											Values: map[string]*SetValue{
												"a": {Expiration: time.Second * 10},
												"b": nil,
												"c": nil,
											},
										},
										"set3": {
											Values: map[string]*SetValue{
												"x": {Expiration: time.Second * 5},
												"y": nil,
											},
										},
									},
								},
								Maps: &Maps{
									Values: map[string]*MapRecordValue{
										"map1": {
											Values: map[string]string{
												"a": "1",
												"b": "2",
											},
										},
										"map2": {
											Values: map[string]string{
												"c": "3",
												"d": "4",
											},
										},
									},
								},
							},
							2: {
								Keys: &Keys{
									Values: map[string]*KeyValue{
										"key3": {
											Value: "value3",
										},
										"key4": {
											Value:      "value4",
											Expiration: time.Second * 5,
										},
									},
								},
								Sets: &Sets{
									Values: map[string]*SetRecordValue{
										"set2": {
											Values: map[string]*SetValue{
												"d": nil,
												"e": {Expiration: time.Second * 5},
												"f": nil,
											},
										},
									},
								},
								Maps: &Maps{
									Values: map[string]*MapRecordValue{
										"map3": {
											Values: map[string]string{
												"c": "3",
												"d": "4",
											},
										},
										"map4": {
											Values: map[string]string{
												"e": "10",
												"f": "11",
											},
										},
									},
								},
							},
						},
					},
				},
				ctx: &context{
					keyRefs: map[string]Keys{},
					setRefs: map[string]SetRecordValue{},
					mapRefs: map[string]MapRecordValue{},
				},
			},
		},
		{
			name: "extend",
			args: args{
				fixtures: []string{"redis_extend"},
			},
			want: want{
				fixtures: []*Fixture{
					{
						Templates: Templates{
							Keys: []*Keys{
								{
									Name:   "parentKeys",
									Values: map[string]*KeyValue{"a": {Value: "1"}, "b": {Value: "2"}},
								},
								{
									Name:   "childKeys",
									Extend: "parentKeys",
									Values: map[string]*KeyValue{"a": {Value: "1"}, "b": {Value: "2"}, "c": {Value: "3"}, "d": {Value: "4"}},
								},
							},
							Sets: []*SetRecordValue{
								{
									Name:   "parentSet",
									Values: map[string]*SetValue{"a": {Expiration: time.Second * 10}, "b": nil},
								},
								{
									Name:   "childSet",
									Extend: "parentSet",
									Values: map[string]*SetValue{"a": {Expiration: time.Second * 10}, "b": nil, "c": nil},
								},
							},
							Maps: []*MapRecordValue{
								{
									Name:   "parentMap",
									Values: map[string]string{"a1": "1", "b1": "2"},
								},
								{
									Name:   "childMap",
									Extend: "parentMap",
									Values: map[string]string{"a1": "1", "b1": "2", "c1": "3"},
								},
							},
						},
						Databases: map[int]Database{
							1: {
								Keys: &Keys{
									Extend: "childKeys",
									Values: map[string]*KeyValue{
										"a": {Value: "1"},
										"b": {Value: "2"},
										"c": {Value: "3"},
										"d": {Value: "4"},
										"key1": {
											Value: "value1",
										},
										"key2": {
											Value:      "value2",
											Expiration: time.Second * 10,
										},
									},
								},
								Sets: &Sets{
									Values: map[string]*SetRecordValue{
										"set1": {
											Extend: "childSet",
											Values: map[string]*SetValue{
												"a": {Expiration: time.Second * 10},
												"b": nil,
												"c": nil,
												"d": nil,
											},
										},
										"set2": {
											Values: map[string]*SetValue{
												"x": nil,
												"y": {Expiration: time.Second * 10},
											},
										},
									},
								},
								Maps: &Maps{
									Values: map[string]*MapRecordValue{
										"map1": {
											Name:   "baseMap",
											Extend: "childMap",
											Values: map[string]string{
												"a":  "1",
												"b":  "2",
												"a1": "1",
												"b1": "2",
												"c1": "3",
											},
										},
										"map2": {
											Values: map[string]string{
												"c": "3",
												"d": "4",
											},
										},
									},
								},
							},
						},
					},
				},
				ctx: &context{
					keyRefs: map[string]Keys{
						"parentKeys": {
							Values: map[string]*KeyValue{
								"a": {Value: "1"},
								"b": {Value: "2"},
							},
						},
						"childKeys": {
							Values: map[string]*KeyValue{
								"a": {Value: "1"},
								"b": {Value: "2"},
								"c": {Value: "3"},
								"d": {Value: "4"},
							},
						},
					},
					setRefs: map[string]SetRecordValue{
						"parentSet": {
							Values: map[string]*SetValue{
								"a": {Expiration: time.Second * 10},
								"b": nil,
							},
						},
						"childSet": {
							Values: map[string]*SetValue{
								"a": {Expiration: time.Second * 10},
								"b": nil,
								"c": nil,
							},
						},
					},
					mapRefs: map[string]MapRecordValue{
						"baseMap": {
							Values: map[string]string{"a": "1", "a1": "1", "b": "2", "b1": "2", "c1": "3"},
						},
						"parentMap": {
							Values: map[string]string{"a1": "1", "b1": "2"},
						},
						"childMap": {
							Values: map[string]string{"a1": "1", "b1": "2", "c1": "3"},
						},
					},
				},
			},
		},
		{
			name: "inherits",
			args: args{
				fixtures: []string{"redis_inherits"},
			},
			want: want{
				fixtures: []*Fixture{
					{
						Inherits: []string{"redis_extend"},
						Databases: map[int]Database{
							1: {
								Keys: &Keys{
									Extend: "childKeys",
									Values: map[string]*KeyValue{
										"a":    {Value: "1"},
										"b":    {Value: "2"},
										"c":    {Value: "3"},
										"d":    {Value: "4"},
										"key1": {Value: "value1"},
										"key2": {Value: "value2", Expiration: time.Second * 10},
									},
								},
								Sets: &Sets{
									Values: map[string]*SetRecordValue{
										"set1": {
											Extend: "childSet",
											Values: map[string]*SetValue{
												"a": {Expiration: time.Second * 10},
												"b": nil,
												"c": nil,
											},
										},
									},
								},
								Maps: &Maps{
									Values: map[string]*MapRecordValue{
										"map1": {
											Extend: "baseMap",
											Values: map[string]string{"a": "1", "a1": "1", "b": "2", "b1": "2", "c1": "3", "x": "10", "y": "11"},
										},
										"map2": {
											Extend: "childMap",
											Values: map[string]string{"a1": "1", "b1": "2", "c1": "3", "j": "1000", "t": "500"},
										},
									},
								},
							},
						},
					},
				},
				ctx: &context{
					keyRefs: map[string]Keys{
						"parentKeys": {
							Values: map[string]*KeyValue{
								"a": {Value: "1"},
								"b": {Value: "2"},
							},
						},
						"childKeys": {
							Values: map[string]*KeyValue{
								"a": {Value: "1"},
								"b": {Value: "2"},
								"c": {Value: "3"},
								"d": {Value: "4"},
							},
						},
					},
					setRefs: map[string]SetRecordValue{
						"parentSet": {
							Values: map[string]*SetValue{
								"a": {Expiration: time.Second * 10},
								"b": nil,
							},
						},
						"childSet": {
							Values: map[string]*SetValue{
								"a": {Expiration: time.Second * 10},
								"b": nil,
								"c": nil,
							},
						},
					},
					mapRefs: map[string]MapRecordValue{
						"baseMap": {
							Values: map[string]string{"a": "1", "a1": "1", "b": "2", "b1": "2", "c1": "3"},
						},
						"parentMap": {
							Values: map[string]string{"a1": "1", "b1": "2"},
						},
						"childMap": {
							Values: map[string]string{"a1": "1", "b1": "2", "c1": "3"},
						},
					},
				},
			},
		},
	}

	p := New([]string{"../../testdata"})

	// test parsing example file from README
	_, err := p.ParseFiles(NewContext(), []string{"redis_example"})
	if err != nil {
		t.Errorf("example file test error: %s", err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := NewContext()
			fixtures, err := p.ParseFiles(ctx, test.args.fixtures)
			if err != nil {
				t.Errorf("ParseFiles - unexpected error: %s", err)
				return
			}
			if diff := cmp.Diff(test.want.fixtures, fixtures); diff != "" {
				t.Errorf("ParseFiles - unexpected diff in fixtures: %s", diff)
			}
			if diff := cmp.Diff(test.want.ctx, ctx, cmp.AllowUnexported(context{})); diff != "" {
				t.Errorf("ParseFiles - unexpected diff in context: %s", diff)
			}
		})
	}
}
