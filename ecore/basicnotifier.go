package ecore

type basicNotifierAdapterList struct {
	*basicEList
	notifier *BasicNotifier
}

type basicNotifierNotification struct {
	*abstractNotification
	notifier *BasicNotifier
}

func (n *basicNotifierNotification) GetNotifier() ENotifier {
	return n.notifier.interfaces.(ENotifier)
}

func (n *basicNotifierNotification) GetFeature() EStructuralFeature {
	return nil
}

func (n *basicNotifierNotification) GetFeatureID() int {
	return -1
}

func newBasicNotifierNotification(notifier *BasicNotifier, eventType EventType, oldValue interface{}, newValue interface{}, position int) *basicNotifierNotification {
	n := new(basicNotifierNotification)
	n.abstractNotification = NewAbstractNotification(eventType, oldValue, newValue, position)
	n.notifier = notifier
	return n
}

func newBasicNotifierAdapterList(notifier *BasicNotifier) *basicNotifierAdapterList {
	l := new(basicNotifierAdapterList)
	l.basicEList = NewEmptyBasicEList()
	l.notifier = notifier
	l.interfaces = l
	return l
}

func (l *basicNotifierAdapterList) didAdd(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	elem.(EAdapter).SetTarget(notifier)
}

func (l *basicNotifierAdapterList) didRemove(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	adapter := elem.(EAdapter)
	if notifier.EDeliver() {
		adapter.NotifyChanged(newBasicNotifierNotification(l.notifier, REMOVING_ADAPTER, elem, nil, index))
	}
	adapter.UnSetTarget(notifier)
}

type BasicNotifier struct {
	interfaces interface{}
	eDeliver   bool
	eAdapters  *basicNotifierAdapterList
}

func NewBasicNotifier() *BasicNotifier {
	notifier := new(BasicNotifier)
	notifier.interfaces = notifier
	notifier.eDeliver = true
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
	if notifier.eAdapters == nil {
		notifier.eAdapters = newBasicNotifierAdapterList(notifier)
	}
	return notifier.eAdapters
}

func (notifier *BasicNotifier) EDeliver() bool {
	return notifier.eDeliver
}

func (notifier *BasicNotifier) ESetDeliver(value bool) {
	notifier.eDeliver = value
}

func (notifier *BasicNotifier) ENotify(notification ENotification) {
	if notifier.eAdapters != nil && notifier.eDeliver {
		for it := notifier.eAdapters.Iterator(); it.HasNext(); {
			it.Next().(EAdapter).NotifyChanged(notification)
		}
	}
}

func (notifier *BasicNotifier) ENotificationRequired() bool {
	return notifier.eAdapters != nil && notifier.eAdapters.Size() > 0 && notifier.eDeliver
}
