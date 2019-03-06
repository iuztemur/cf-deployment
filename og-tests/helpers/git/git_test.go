package git_test

import (
  "errors"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/onsi/gomega/gstruct"
  "gopkg.in/src-d/go-git.v4/plumbing"
  "gopkg.in/src-d/go-git.v4/plumbing/storer"
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
		var (
			expectedVersion string
			expectedErr     error
		)

		JustBeforeEach(func() {
      reference := plumbing.NewReferenceFromStrings("DUMMY PLACEHOLDER REFERENCE","refs/tags/")
			expectedVersion, expectedErr = helper.GetMajorVersion(reference)
		})

		Context("when git returns a list of valid tag", func() {
			BeforeEach(func() {
				fake := new(gitfakes.FakeRepository)

				reference1 := plumbing.NewReferenceFromStrings("v1.0.0","refs/tags/")
        reference2 := plumbing.NewReferenceFromStrings("v1.0.1","refs/tags/")

        tags := storer.NewReferenceSliceIter([]*plumbing.Reference{reference1, reference2})

				fake.TagsReturns(tags, nil)

				//client = fake
			})

			It("returns the most recent major version", func() {
				Expect(expectedErr).To(BeNil())
				Expect(expectedVersion).To(Equal("v1.0.0"))
			})
		})

		Context("when git returns an invalid tag", func() {
      BeforeEach(func() {
        fake := new(gitfakes.FakeRepository)

        reference1 := plumbing.NewReferenceFromStrings("invalid","refs/tags/")
        tags := storer.NewReferenceSliceIter([]*plumbing.Reference{reference1})

        fake.TagsReturns(tags, nil)

        //client = fake
      })

      It("returns the most recent major version", func() {
        Expect(expectedErr).To(MatchError("invalid version"))
        Expect(expectedVersion).To(nil)
      })

    })
  //
	//	Context("when git returns an error", func() {
	//		BeforeEach(func() {
  //      fake := new(gitfakes.FakeRepository)
  //      fake.TagsReturns(nil, errors.New("some error"))
  //
	//			//client = fake
	//		})
  //
	//		It("returns an error", func() {
	//			Expect(expectedErr).To(MatchError("some error"))
	//		})
	//	})
	})

	//Describe("CheckoutSubdir", func() {})
})
