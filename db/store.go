package db

import (
	"github.com/antaresvision/helloserver/models"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func (ds *Store) Close() error {
	return ds.db.Close()
}

func (ds *Store) CreateItem(ntinid int, serial string, status int) (models.Item, error) {
	item := models.Item{
		NtinId: ntinid,
		Serial: serial,
		Status: status,
	}

	nstmt, err := ds.db.PrepareNamed(`INSERT INTO public.items(ntinid,serial,status) 
											VALUES(:ntinid, :serial, :status) RETURNING id`)
	if err != nil {
		return models.Item{}, err
	}

	err = nstmt.QueryRowx(item).Scan(&item.Id)
	return item, err
}

func (ds *Store) GetItem(id int) (models.Item, error) {
	item := models.Item{
		Id:id,
	}

	nstmt, err := ds.db.PrepareNamed(`SELECT * FROM public.items WHERE id=:id`)
	if err != nil {
		return models.Item{}, err
	}

	err = nstmt.Get(&item, item)
	return item, err
}


func (ds *Store) GetAllItems() ([]models.Item, error) {
	items := []models.Item{}
	err := ds.db.Select(&items, `SELECT * FROM public.items`)
	return items, err
}

func (ds *Store) UpdateItem(item models.Item) error {
	_, err := ds.db.NamedExec(`UPDATE items SET ntinid=:ntinid, serial=:serial, status=:status WHERE id=:id`, item)
	return err
}

func (ds *Store) RemoveItem(id int) error {
	item := models.Item{
		Id: id,
	}
	_, err := ds.db.NamedExec(`DELETE FROM items WHERE id=:id`, item)
	return err
}