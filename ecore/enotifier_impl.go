package ecore

type ENotifierImpl struct {
	AbstractENotifier
	deliver  bool
	adapters *notifierAdapterList
}

func (notifier *ENotifierImpl) Initialize() {
	notifier.deliver = true
	notifier.adapters = nil
}

func (notifier *ENotifierImpl) EAdapters() EList {
	if notifier.adapters == nil {
		notifier.adapters = newNotifierAdapterList(&notifier.AbstractENotifier)
	}
	return notifier.adapters
}

func (notifier *ENotifierImpl) EBasicAdapters() EList {
	if notifier.adapters == nil {
		return nil
	}
	return notifier.adapters
}

func (notifier *ENotifierImpl) EDeliver() bool {
	return notifier.deliver
}

func (notifier *ENotifierImpl) ESetDeliver(deliver bool) {
	notifier.deliver = deliver
}
