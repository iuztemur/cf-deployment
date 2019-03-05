package git_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "og/helpers/git"
	"og/helpers/git/gitfakes"
)

var _ = Describe("Git", func() {
	var (
		repo Repository

		helper *Helper
	)

	JustBeforeEach(func() {
		helper = NewHelper(repo)
	})

	Describe("SetupRepo", func() {
		var (
			fake *gitfakes.FakeRepository

			expectedErr error
		)

		BeforeEach(func() {
			fake = new(gitfakes.FakeRepository)

			repo = fake
		})

		JustBeforeEach(func() {
			expectedErr = helper.SetupRepository()
		})

		It("attempts to add cf-deployment as the origin remote", func() {
			Expect(fake.CreateRemoteCallCount()).To(Equal(1))

			config := fake.CreateRemoteArgsForCall(0)
			Expect(*config).To(MatchFields(IgnoreExtras, Fields{
				"Name": Equal("origin"),
				"URLs": ConsistOf("https://github.com/cloudfoundry/cf-deployment.git"),
			}))
		})

		Context("when adding the remote is successful", func() {
			BeforeEach(func() {
				fake.CreateRemoteReturns(nil, nil)
			})

			It("returns successfully", func() {
				Expect(expectedErr).To(BeNil())
			})
		})

		Context("when adding the remote fails", func() {
			BeforeEach(func() {
				fake.CreateRemoteReturns(nil, errors.New("some git error"))
			})

			It("it returns an error", func() {
				Expect(expectedErr).To(MatchError("error adding remote: some git error"))
			})
		})
	})

	Describe("GetMajorVersion", func() {
		// var (
		// 	expectedVersion string
		// 	expectedErr     error
		// )

		JustBeforeEach(func() {
			// expectedVersion, expectedErr = helper.GetMajorVersion()
		})

		// Context("when git returns a list of valid tag", func() {
		// 	BeforeEach(func() {
		// 		fake := new(gitfakes.FakeClient)
		// 		fake.FetchContextReturns(nil)

		// 		client = fake
		// 	})

		// 	It("returns the most recent major version", func() {
		// 		Expect(expectedErr).To(BeNil())
		// 		Expect(expectedVersion).To(Equal("v1.0.0"))
		// 	})
		// })

		// Context("when git returns an invalid tag", func() {})

		// Context("when git returns an error", func() {
		// 	BeforeEach(func() {
		// 		fake := new(gitfakes.FakeClient)
		// 		fake.FetchContextReturns(errors.New("some error"))

		// 		client = fake
		// 	})

		// 	It("returns an error", func() {
		// 		Expect(expectedErr).To(MatchError("some error"))
		// 	})
		// })
	})

	Describe("CheckoutSubdir", func() {})
})
