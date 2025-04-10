package api

import (
	"reflect"
	"testing"
)

func Test_parseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want *fieldNode
	}{
		{
			name: "Test 1: id,name",
			args: args{query: "id,name"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "name": &fieldNode{
					Name:      "name",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 2: id,name{address{city,zip},orders{id,total}}",
			args: args{query: "id,name{address{city,zip},orders{id,total}}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "name": &fieldNode{
					Name: "name",
					Children: map[string]*fieldNode{"address": &fieldNode{
						Name: "address",
						Children: map[string]*fieldNode{"city": &fieldNode{
							Name:      "city",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "zip": &fieldNode{
							Name:      "zip",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}, "orders": &fieldNode{
						Name: "orders",
						Children: map[string]*fieldNode{"id": &fieldNode{
							Name:      "id",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "total": &fieldNode{
							Name:      "total",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 3: id,name,address{city,zip},orders{id,total}",
			args: args{query: "id,name,address{city,zip},orders{id,total}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "name": &fieldNode{
					Name:      "name",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "address": &fieldNode{
					Name: "address",
					Children: map[string]*fieldNode{"city": &fieldNode{
						Name:      "city",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "zip": &fieldNode{
						Name:      "zip",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}},
					Requested: true,
				}, "orders": &fieldNode{
					Name: "orders",
					Children: map[string]*fieldNode{"id": &fieldNode{
						Name:      "id",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "total": &fieldNode{
						Name:      "total",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 4: id,name,details{specs{width,height},color}",
			args: args{query: "id,name,details{specs{width,height},color}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "name": &fieldNode{
					Name:      "name",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "details": &fieldNode{
					Name: "details",
					Children: map[string]*fieldNode{"specs": &fieldNode{
						Name: "specs",
						Children: map[string]*fieldNode{"width": &fieldNode{
							Name:      "width",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "height": &fieldNode{
							Name:      "height",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}, "color": &fieldNode{
						Name:      "color",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 5: user{id,profile{name,email},posts{id,title,comments{id,text}}}",
			args: args{query: "user{id,profile{name,email},posts{id,title,comments{id,text}}}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"user": &fieldNode{
					Name: "user",
					Children: map[string]*fieldNode{"id": &fieldNode{
						Name:      "id",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "profile": &fieldNode{
						Name: "profile",
						Children: map[string]*fieldNode{"name": &fieldNode{
							Name:      "name",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "email": &fieldNode{
							Name:      "email",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}, "posts": &fieldNode{
						Name: "posts",
						Children: map[string]*fieldNode{"id": &fieldNode{
							Name:      "id",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "title": &fieldNode{
							Name:      "title",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "comments": &fieldNode{
							Name: "comments",
							Children: map[string]*fieldNode{"id": &fieldNode{
								Name:      "id",
								Children:  map[string]*fieldNode{},
								Requested: true,
							}, "text": &fieldNode{
								Name:      "text",
								Children:  map[string]*fieldNode{},
								Requested: true,
							}},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 6: simple",
			args: args{query: "simple"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"simple": &fieldNode{
					Name:      "simple",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 7: nested{a,b{c,d}}",
			args: args{query: "nested{a,b{c,d}}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"nested": &fieldNode{
					Name: "nested",
					Children: map[string]*fieldNode{"a": &fieldNode{
						Name:      "a",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "b": &fieldNode{
						Name: "b",
						Children: map[string]*fieldNode{"c": &fieldNode{
							Name:      "c",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "d": &fieldNode{
							Name:      "d",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 8: a,b,c",
			args: args{query: "a,b,c"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"a": &fieldNode{
					Name:      "a",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "b": &fieldNode{
					Name:      "b",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "c": &fieldNode{
					Name:      "c",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 4: id,product{name,category{title},brand{title,country}},profile{name,email},posts{id,title,comments{id,text}}",
			args: args{query: "id,product{name,category{title},brand{title,country}},profile{name,email},posts{id,title,comments{id,text}}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "product": &fieldNode{
					Name: "product",
					Children: map[string]*fieldNode{"name": &fieldNode{
						Name:      "name",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "category": &fieldNode{
						Name: "category",
						Children: map[string]*fieldNode{"title": &fieldNode{
							Name:      "title",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}, "brand": &fieldNode{
						Name: "brand",
						Children: map[string]*fieldNode{"title": &fieldNode{
							Name:      "title",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "country": &fieldNode{
							Name:      "country",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}, "profile": &fieldNode{
					Name: "profile",
					Children: map[string]*fieldNode{"name": &fieldNode{
						Name:      "name",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "email": &fieldNode{
						Name:      "email",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}},
					Requested: true,
				}, "posts": &fieldNode{
					Name: "posts",
					Children: map[string]*fieldNode{"id": &fieldNode{
						Name:      "id",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "title": &fieldNode{
						Name:      "title",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "comments": &fieldNode{
						Name: "comments",
						Children: map[string]*fieldNode{"id": &fieldNode{
							Name:      "id",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "text": &fieldNode{
							Name:      "text",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
		{name: "Test 5: id,product{name,category{title,ref{code{date,id,isbn{id,date,url{date,ref}}}}}}",
			args: args{query: "id,product{name,category{title,ref{code{date,id,isbn{id,date,url{date,ref}}}}}}"},
			want: &fieldNode{
				Name: "root",
				Children: map[string]*fieldNode{"id": &fieldNode{
					Name:      "id",
					Children:  map[string]*fieldNode{},
					Requested: true,
				}, "product": &fieldNode{
					Name: "product",
					Children: map[string]*fieldNode{"name": &fieldNode{
						Name:      "name",
						Children:  map[string]*fieldNode{},
						Requested: true,
					}, "category": &fieldNode{
						Name: "category",
						Children: map[string]*fieldNode{"title": &fieldNode{
							Name:      "title",
							Children:  map[string]*fieldNode{},
							Requested: true,
						}, "ref": &fieldNode{
							Name: "ref",
							Children: map[string]*fieldNode{"code": &fieldNode{
								Name: "code",
								Children: map[string]*fieldNode{"date": &fieldNode{
									Name:      "date",
									Children:  map[string]*fieldNode{},
									Requested: true,
								}, "id": &fieldNode{
									Name:      "id",
									Children:  map[string]*fieldNode{},
									Requested: true,
								}, "isbn": &fieldNode{
									Name: "isbn",
									Children: map[string]*fieldNode{"id": &fieldNode{
										Name:      "id",
										Children:  map[string]*fieldNode{},
										Requested: true,
									}, "date": &fieldNode{
										Name:      "date",
										Children:  map[string]*fieldNode{},
										Requested: true,
									}, "url": &fieldNode{
										Name: "url",
										Children: map[string]*fieldNode{"date": &fieldNode{
											Name:      "date",
											Children:  map[string]*fieldNode{},
											Requested: true,
										}, "ref": &fieldNode{
											Name:      "ref",
											Children:  map[string]*fieldNode{},
											Requested: true,
										}},
										Requested: true,
									}},
									Requested: true,
								}},
								Requested: true,
							}},
							Requested: true,
						}},
						Requested: true,
					}},
					Requested: true,
				}},
				Requested: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseQuery(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
