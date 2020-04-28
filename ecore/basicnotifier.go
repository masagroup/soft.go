package ecore

type adapterList struct {
	*basicEList
	notifier *BasicNotifier
}

type adapterNotification struct {
	*abstractNotification
	notifier ENotifier
}

func (n *adapterNotification) GetNotifier() ENotifier {
	return n.notifier
}

func (n *adapterNotification) GetFeature() EStructuralFeature {
	return nil
}

func (n *adapterNotification) GetFeatureID() int {
	return -1
}

func newAdapterNotification(notifier ENotifier, eventType EventType, oldValue interface{}, newValue interface{}, position int) *adapterNotification {
	n := new(adapterNotification)
	n.abstractNotification = NewAbstractNotification(eventType, oldValue, newValue, position)
	n.notifier = notifier
	return n
}

func newAdapterList(notifier *BasicNotifier) *adapterList {
	l := new(adapterList)
	l.basicEList = NewEmptyBasicEList()
	l.notifier = notifier
	l.interfaces = l
	return l
}

func (l *adapterList) didAdd(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	elem.(EAdapter).SetTarget(notifier)
}

func (l *adapterList) didRemove(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	adapter := elem.(EAdapter)
	if notifier.EDeliver() {
		adapter.NotifyChanged(newAdapterNotification(notifier, REMOVING_ADAPTER, elem, nil, index))
	}
	adapter.UnSetTarget(notifier)
}

type BasicNotifier struct {
	interfaces interface{}
	eDeliver   bool
	eAdapters  EList
}

func NewBasicNotifier() *BasicNotifier {
	notifier := new(BasicNotifier)
	notifier.interfaces = notifier
	notifier.eDeliver = true
	notifier.eAdapters = newAdapterList(notifier)
	return notifier
}

// SetInterfaces ...
func (o *BasicNotifier) SetInterfaces(interfaces interface{}) {
	o.interfaces = interfaces
}

// GetInterfaces ...
func (o *BasicNotifier) GetInterfaces() interface{} {
	return o.interfaces
}

func (notifier *BasicNotifier) EAdapters() EList {
	return notifier.eAdapters
}

func (notifier *BasicNotifier) EDeliver() bool {
	return notifier.eDeliver
}

func (notifier *BasicNotifier) ESetDeliver(value bool) {
	notifier.eDeliver = value
}

func (notifier *BasicNotifier) ENotify(notification ENotification) {
	for it := notifier.eAdapters.Iterator(); it.HasNext(); {
		it.Next().(EAdapter).NotifyChanged(notification)
	}
}

func (notifier *BasicNotifier) ENotificationRequired() bool {
	return notifier.eAdapters != nil && notifier.eAdapters.Size() > 0
}
