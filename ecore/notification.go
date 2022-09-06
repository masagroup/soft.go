package ecore

type notification struct {
	AbstractNotification
	object    EObject
	feature   EStructuralFeature
	featureID int
}

// NewNotificationByFeature ...
func NewNotificationByFeature(object EObject, eventType EventType, feature EStructuralFeature, oldValue any, newValue any, position int) *notification {
	n := new(notification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.object = object
	n.feature = feature
	n.featureID = -1
	return n
}

// NewNotificationByFeatureID ...
func NewNotificationByFeatureID(object EObject, eventType EventType, featureID int, oldValue any, newValue any, position int) *notification {
	n := new(notification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.object = object
	n.feature = nil
	n.featureID = featureID
	return n
}

func (notif *notification) GetNotifier() ENotifier {
	return notif.object.(ENotifier)
}

func (notif *notification) GetFeature() EStructuralFeature {
	if notif.feature != nil {
		return notif.feature
	}
	return notif.object.EClass().GetEStructuralFeature(notif.featureID)
}

func (notif *notification) GetFeatureID() int {
	if notif.featureID != -1 {
		return notif.featureID
	}
	if notif.feature != nil {
		return notif.feature.GetFeatureID()
	}
	return -1
}
