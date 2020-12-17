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
}

type ENotifierInternal interface {
	ENotifier
	HasAdapters() bool
}

func NewBasicNotifier() *BasicNotifier {
	notifier := new(BasicNotifier)
	notifier.interfaces = notifier
	return notifier
}

func (notifier *BasicNotifier) AsENotifier() ENotifier {
	return notifier.interfaces.(ENotifier)
}

// SetInterfaces ...
func (notifier *BasicNotifier) SetInterfaces(interfaces interface{}) {
	notifier.interfaces = interfaces
}

// GetInterfaces ...
func (notifier *BasicNotifier) GetInterfaces() interface{} {
	return notifier.interfaces
}

func (notifier *BasicNotifier) HasAdapters() bool {
	adapters := notifier.AsENotifier().EAdapters()
	return adapters != nil && !adapters.Empty()
}

func (notifier *BasicNotifier) EAdapters() EList {
	return NewEmptyImmutableEList()
}

func (notifier *BasicNotifier) EDeliver() bool {
	return false
}

func (notifier *BasicNotifier) ESetDeliver(value bool) {
	panic("operation not supported")
}

func (notifier *BasicNotifier) ENotify(notification ENotification) {
	n := notifier.AsENotifier()
	adapters := n.EAdapters()
	deliver := n.EDeliver()
	if adapters != nil && deliver {
		for it := adapters.Iterator(); it.HasNext(); {
			it.Next().(EAdapter).NotifyChanged(notification)
		}
	}
}

func (notifier *BasicNotifier) ENotificationRequired() bool {
	n := notifier.interfaces.(ENotifierInternal)
	return n.HasAdapters() && n.EDeliver()
}
