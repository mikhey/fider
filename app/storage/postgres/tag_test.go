package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTagStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	tag, err := tags.Add("Feature Request", "FF0000", true)
	Expect(err).To(BeNil())
	Expect(tag.ID).To(Equal(1))

	dbTag, err := tags.GetBySlug("feature-request")

	Expect(err).To(BeNil())
	Expect(dbTag.ID).To(Equal(1))
	Expect(dbTag.Name).To(Equal("Feature Request"))
	Expect(dbTag.Slug).To(Equal("feature-request"))
	Expect(dbTag.Color).To(Equal("FF0000"))
	Expect(dbTag.IsPublic).To(BeTrue())
}

func TestTagStorage_AddUpdateAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	tag, err := tags.Add("Feature Request", "FF0000", true)
	tag, err = tags.Update(tag.ID, "Bug", "000000", false)

	dbTag, err := tags.GetBySlug("bug")

	Expect(err).To(BeNil())
	Expect(dbTag.ID).To(Equal(tag.ID))
	Expect(dbTag.Name).To(Equal("Bug"))
	Expect(dbTag.Slug).To(Equal("bug"))
	Expect(dbTag.Color).To(Equal("000000"))
	Expect(dbTag.IsPublic).To(BeFalse())
}

func TestTagStorage_AddRemoveAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	tag, err := tags.Add("Bug", "FFFFFF", true)

	err = tags.Remove(tag.ID)
	Expect(err).To(BeNil())

	dbTag, err := tags.GetBySlug("bug")

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbTag).To(BeNil())
}

func TestTagStorage_Assign_Unassign(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	ideas := postgres.NewIdeaStorage(demoTenant(tenants), trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	idea, _ := ideas.Add("My great idea", "with a great description", 2)
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag.ID, idea.ID, 2)
	Expect(err).To(BeNil())

	assigned, err := tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(1))
	Expect(assigned[0].ID).To(Equal(tag.ID))
	Expect(assigned[0].Name).To(Equal("Bug"))
	Expect(assigned[0].Slug).To(Equal("bug"))
	Expect(assigned[0].Color).To(Equal("FFFFFF"))
	Expect(assigned[0].IsPublic).To(BeTrue())

	err = tags.UnassignTag(tag.ID, idea.ID)
	Expect(err).To(BeNil())

	assigned, err = tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(0))
}

func TestTagStorage_Assign_RemoveTag(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	ideas := postgres.NewIdeaStorage(demoTenant(tenants), trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	idea, _ := ideas.Add("My great idea", "with a great description", 2)
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag.ID, idea.ID, 2)
	Expect(err).To(BeNil())

	err = tags.Remove(tag.ID)
	Expect(err).To(BeNil())

	assigned, err := tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(0))
}
