package router_test

import (
	. "Gateway311/router"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {

	Context("Valid JSON / 1 Provider / 3 Services", func() {
		var file1 string = `
		{
			"serviceCategories": [
				"Abandoned",
				"Eyesore",
				"Graffiti",
				"Street",
				"Trash"
			],
			"serviceAreas": {
				"san jose": {
					"id": 1,
					"name": "San Jose",
					"serviceProviders": {
						"citysourced": {
							"id": 1,
							"name": "CitySourced - SJ",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 1,
								"name": "Abandoned Bicycle",
								"catg": ["Abandoned"]
							}, {
								"id": 2,
								"name": "Abandoned Car",
								"catg": ["Abandoned"]
							}, {
								"id": 3,
								"name": "Abandoned Home",
								"catg": ["Abandoned"]
							}]
						}
					}
				}
			}
		}
		`
		var rd RouteData

		It("should load OK", func() {
			err := rd.Load([]byte(file1))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have 3 services for San Jose", func() {
			id, services, err := rd.Services("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(1))
			Expect(services).To(HaveLen(3))
		})

		It("should have 1 Service Provider for San Jose", func() {
			providers, err := rd.ServiceProviders("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(providers).To(HaveLen(1))
			Expect(providers[0].Name).To(Equal("CitySourced - SJ"))
		})

		It("should have 0 services for Morgan Hill", func() {
			id, services, err := rd.Services("Morgan Hill")
			Expect(err).To(HaveOccurred())
			Expect(id).To(Equal(0))
			Expect(services).To(HaveLen(0))
		})

		It("should have 0 Service Providers for Morgan Hill", func() {
			providers, err := rd.ServiceProviders("Morgan Hill")
			Expect(err).To(HaveOccurred())
			Expect(providers).To(HaveLen(0))
		})

		It("should have have CitySourced provider for Service 2", func() {
			provider, err := rd.ServiceProvider(2)
			Expect(err).NotTo(HaveOccurred())
			Expect(provider.ID).To(Equal(1))
			Expect(provider.Name).To(Equal("CitySourced - SJ"))
		})

		It("should have have CitySourced provider interface type for Service 2", func() {
			itype, err := rd.ServiceProviderInterface(2)
			Expect(err).NotTo(HaveOccurred())
			Expect(itype).To(Equal("CitySourced"))
		})
	})

	Context("Valid JSON / 2 Providers / 8 Services / San Jose", func() {
		var file1 string = `
		{
			"serviceCategories": [
				"Abandoned",
				"Eyesore",
				"Graffiti",
				"Street",
				"Trash"
			],
			"serviceAreas": {
				"san jose": {
					"id": 1,
					"name": "San Jose",
					"serviceProviders": {
						"citysourced": {
							"id": 1,
							"name": "CitySourced - SJ",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 1,
								"name": "Abandoned Bicycle",
								"catg": ["Abandoned"]
							}, {
								"id": 2,
								"name": "Abandoned Car",
								"catg": ["Abandoned"]
							}, {
								"id": 3,
								"name": "Abandoned Home",
								"catg": ["Abandoned"]
							}]
						},
		        		"citysourced2": {
							"id": 2,
							"name": "CitySourced2 - SJ",
							"interfaceType": "SeeClickFix",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 36,
								"name": "Information Only (Please Describe)",
								"catg": []
							}, {
								"id": 37,
								"name": "Other (Not Listed Please Describe)",
								"catg": []
							}, {
								"id": 38,
								"name": "Other (Please Describe)",
								"catg": []
							}, {
								"id": 39,
								"name": "Parking Issue",
								"catg": []
							}, {
								"id": 40,
								"name": "Parking Issues",
								"catg": []
							}]
						}
					}
				}
			}
		}
		`
		var rd RouteData

		It("should load OK", func() {
			err := rd.Load([]byte(file1))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have 8 services for San Jose", func() {
			id, services, err := rd.Services("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(1))
			Expect(services).To(HaveLen(8))
		})

		It("should have 2 Service Provider for San Jose", func() {
			providers, err := rd.ServiceProviders("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(providers).To(HaveLen(2))
			Expect(providers[0].Name).To(Equal("CitySourced - SJ"))
			Expect(providers[1].Name).To(Equal("CitySourced2 - SJ"))
		})

		It("should have have CitySourced provider interface type for Service 2", func() {
			itype, err := rd.ServiceProviderInterface(2)
			Expect(err).NotTo(HaveOccurred())
			Expect(itype).To(Equal("CitySourced"))
		})
		It("should have have SeeClickFix provider interface type for Service 37", func() {
			itype, err := rd.ServiceProviderInterface(37)
			Expect(err).NotTo(HaveOccurred())
			Expect(itype).To(Equal("SeeClickFix"))
		})
		It("should have have unknow provider interface type for Service 10", func() {
			itype, err := rd.ServiceProviderInterface(10)
			Expect(err).To(HaveOccurred())
			Expect(itype).To(Equal(""))
		})

	})

	Context("Valid JSON / 2 Areas / 3 Providers / 8 Services in San Jose, 3 Services in Morgan Hil", func() {
		var file1 string = `
		{
			"serviceCategories": [
				"Abandoned",
				"Eyesore",
				"Graffiti",
				"Street",
				"Trash"
			],
			"serviceAreas": {
				"san jose": {
					"id": 1,
					"name": "San Jose",
					"serviceProviders": {
						"citysourced": {
							"id": 1,
							"name": "CitySourced - SJ",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 1,
								"name": "Abandoned Bicycle",
								"catg": ["Abandoned"]
							}, {
								"id": 2,
								"name": "Abandoned Car",
								"catg": ["Abandoned"]
							}, {
								"id": 3,
								"name": "Abandoned Home",
								"catg": ["Abandoned"]
							}]
						},
						"citysourced2": {
							"id": 2,
							"name": "CitySourced2 - SJ",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 36,
								"name": "Information Only (Please Describe)",
								"catg": []
							}, {
								"id": 37,
								"name": "Other (Not Listed Please Describe)",
								"catg": []
							}, {
								"id": 38,
								"name": "Other (Please Describe)",
								"catg": []
							}, {
								"id": 39,
								"name": "Parking Issue",
								"catg": []
							}, {
								"id": 40,
								"name": "Parking Issues",
								"catg": []
							}]
						}
					}
				},
				"morgan hill": {
					"id": 2,
					"name": "Morgan Hill",
					"serviceProviders": {
						"citysourced": {
							"id": 3,
							"name": "CitySourced - MH",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 1,
								"name": "Abandoned Bicycle",
								"catg": ["Abandoned"]
							}, {
								"id": 2,
								"name": "Abandoned Car",
								"catg": ["Abandoned"]
							}, {
								"id": 3,
								"name": "Abandoned Home",
								"catg": ["Abandoned"]
							}]
						}
					}
				}
			}
		}
		`
		var rd RouteData

		It("should load OK", func() {
			err := rd.Load([]byte(file1))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have 8 services for San Jose", func() {
			id, services, err := rd.Services("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(1))
			Expect(services).To(HaveLen(8))
		})

		It("should have 2 Service Provider for San Jose", func() {
			providers, err := rd.ServiceProviders("San Jose")
			Expect(err).NotTo(HaveOccurred())
			Expect(providers).To(HaveLen(2))
			Expect(providers[0].Name).To(Equal("CitySourced - SJ"))
			Expect(providers[1].Name).To(Equal("CitySourced2 - SJ"))
		})

		It("should have 3 services for Morgan Hill", func() {
			id, services, err := rd.Services("Morgan Hill")
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(2))
			Expect(services).To(HaveLen(3))
		})

		It("should have 1 Service Provider for Morgan Hill", func() {
			providers, err := rd.ServiceProviders("Morgan Hill")
			Expect(err).NotTo(HaveOccurred())
			Expect(providers).To(HaveLen(1))
			Expect(providers[0].Name).To(Equal("CitySourced - MH"))
		})
	})

	Context("Invalid JSON", func() {
		var file1 string = `
		{
			"serviceCategories": [
				"Abandoned",
				"Eyesore",
				"Graffiti",
				"Street",
				"Trash"
			],
			"serviceAreas": {
				"san jose": {
					"id": 1x,
					"name": "San Jose",
					"serviceProviders": {
						"citysourced": {
							"id": 1,
							"name": "CitySourced - SJ",
							"interfaceType": "CitySourced",
							"url": "http://localhost:5050/api/",
							"key": "a01234567890z",
							"services": [{
								"id": 1,
								"name": "Abandoned Bicycle",
								"catg": ["Abandoned"]
							}, {
								"id": 2,
								"name": "Abandoned Car",
								"catg": ["Abandoned"]
							}, {
								"id": 3,
								"name": "Abandoned Home",
								"catg": ["Abandoned"]
							}]
						}
					}
				}
			}
		}
		`
		var rd RouteData

		It("should NOT load OK", func() {
			err := rd.Load([]byte(file1))
			Expect(err).To(HaveOccurred())
		})
	})
})
