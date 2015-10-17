package models

func Find(limit, start int, objs interface{}) error {
	return x.Limit(limit, start).Find(objs)
}
func GetById(id int64, obj interface{}) error {
	has, err := x.Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}

	return nil
}
func Get(obj interface{}) error {
	has, err := x.Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func Count(obj interface{}) (int64, error) {
	return x.Count(obj)
}

func IsExist(obj interface{}) bool {
	has, _ := x.Get(obj)
	return has
}

func Insert(obj interface{}) error {
	_, err := x.Insert(obj)
	return err
}

func DeleteById(id int64, obj interface{}) error {
	_, err := x.Id(id).Delete(obj)
	return err
}

func Delete(obj interface{}) error {
	_, err := x.Delete(obj)
	return err
}

func UpdateById(id int64, obj interface{}, cols ...string) error {
	_, err := x.Cols(cols...).Id(id).Update(obj)
	return err

}
