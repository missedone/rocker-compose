package compose

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestComparatorSameValue(t *testing.T) {
	cmp := NewDiff()
	containers := make([]*Container, 0)
	act, err := cmp.Diff("", containers, containers)
	assert.Empty(t, act)
	assert.Nil(t, err)
}

func TestDiffCreateAll(t *testing.T) {
	cmp := NewDiff()
	containers := []*Container{}
	c1 := newContainer("test", "1", ContainerName{"test", "2"}, ContainerName{"test", "3"})
	c2 := newContainer("test", "2", ContainerName{"test", "4"})
	c3 := newContainer("test", "3", ContainerName{"test", "4"})
	c4 := newContainer("test", "4")
	containers = append(containers, c1, c2, c3, c4)
	actions, _ := cmp.Diff("test", containers, []*Container{})
	mock := clientMock{}
	mock.On("CreateContainer", c4).Return(nil)
	mock.On("CreateContainer", c2).Return(nil)
	mock.On("CreateContainer", c3).Return(nil)
	mock.On("CreateContainer", c1).Return(nil)
	runner := NewDockerClientRunner(&mock)
	runner.Run(actions)
	mock.AssertExpectations(t)
}


func TestDiffNoDependencies(t *testing.T) {
	cmp := NewDiff()
	containers := []*Container{}
	c1 := newContainer("test", "1")
	c2 := newContainer("test", "2")
	c3 := newContainer("test", "3")
	containers = append(containers, c1, c2, c3)
	actions, _ := cmp.Diff("test", containers, []*Container{})
	mock := clientMock{}
	mock.On("CreateContainer", c1).Return(nil)
	mock.On("CreateContainer", c2).Return(nil)
	mock.On("CreateContainer", c3).Return(nil)
	runner := NewDockerClientRunner(&mock)
	runner.Run(actions)
	mock.AssertExpectations(t)
}

func TestDiffCreateRemoving(t *testing.T) {
	cmp := NewDiff()
	containers := []*Container{}
	c1 := newContainer("test", "1", ContainerName{"test", "2"}, ContainerName{"test", "3"})
	c2 := newContainer("test", "2", ContainerName{"test", "4"})
	c3 := newContainer("test", "3", ContainerName{"test", "4"})
	c4 := newContainer("test", "4")
	c5 := newContainer("test", "5")
	containers = append(containers, c1, c2, c3, c4)
	actions, _ := cmp.Diff("test", containers, []*Container{c5})
	mock := clientMock{}
	mock.On("RemoveContainer", c5).Return(nil)
	mock.On("CreateContainer", c4).Return(nil)
	mock.On("CreateContainer", c2).Return(nil)
	mock.On("CreateContainer", c3).Return(nil)
	mock.On("CreateContainer", c1).Return(nil)
	runner := NewDockerClientRunner(&mock)
	runner.Run(actions)
	mock.AssertExpectations(t)
}

func TestDiffCreateSome(t *testing.T) {
	cmp := NewDiff()
	containers := []*Container{}
	c1 := newContainer("test", "1", ContainerName{"test", "2"}, ContainerName{"test", "3"})
	c2 := newContainer("test", "2", ContainerName{"test", "4"})
	c3 := newContainer("test", "3", ContainerName{"test", "4"})
	c4 := newContainer("test", "4")
	containers = append(containers, c1, c2, c3, c4)
	actions, _ := cmp.Diff("test", containers, []*Container{c1})
	mock := clientMock{}
	mock.On("CreateContainer", c4).Return(nil)
	mock.On("CreateContainer", c2).Return(nil)
	mock.On("CreateContainer", c3).Return(nil)
	runner := NewDockerClientRunner(&mock)
	runner.Run(actions)
	mock.AssertExpectations(t)
}

func newContainer(namespace string, name string, dependencies ...ContainerName) *Container {
	return &Container{
		State: &ContainerState{
			Running: true,
		},
		Name: &ContainerName{namespace, name},
		Config: &ConfigContainer{
			VolumesFrom: dependencies,
		}}
}

func (m *clientMock) GetContainers() ([]*Container, error) {
	args := m.Called()
	return nil, args.Error(0)
}

func (m *clientMock) RemoveContainer(container *Container) error {
	args := m.Called(container)
	return args.Error(0)
}

func (m *clientMock) CreateContainer(container *Container) error {
	args := m.Called(container)
	return args.Error(0)
}

func (m *clientMock) EnsureContainer(container *Container) error {
	args := m.Called(container)
	return args.Error(0)
}

func (m *clientMock) PullImage(imageName *ImageName) error {
	args := m.Called(imageName)
	return args.Error(0)
}

func (m *clientMock) PullAll(config *Config) error {
	args := m.Called(config)
	return args.Error(0)
}

type clientMock struct {
	mock.Mock
}