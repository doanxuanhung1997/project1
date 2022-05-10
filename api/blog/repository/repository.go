package repository

type masterInterface interface {
	GetAllConsultingField() ()
}

func PublishInterfaceMaster() masterInterface {
	return &masterResource{}
}

type masterResource struct {
}

func (r *masterResource) GetAllConsultingField() () {
}
