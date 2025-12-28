package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	elastic "github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	// No explicit close method for elastic.Client
}

func (r *elasticRepository) PutProduct(ctx context.Context, p *Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	var doc productDocument

	if err := json.Unmarshal(res.Source, &doc); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        doc.Name,
		Description: doc.Description,
		Price:       doc.Price,
	}, err

}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var doc productDocument
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}
		products = append(products, &Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		})
	}

	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]*Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}

	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0, len(res.Docs))
	for _, doc := range res.Docs {
		if !doc.Found {
			continue
		}
		var pd productDocument
		if err := json.Unmarshal(doc.Source, &pd); err != nil {
			return nil, err
		}
		products = append(products, &Product{
			ID:          doc.Id,
			Name:        pd.Name,
			Description: pd.Description,
			Price:       pd.Price,
		})
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)

	if err != nil {
		log.Println("SearchProducts error:", err)
		return nil, err
	}

	products := make([]*Product, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var doc productDocument
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}
		products = append(products, &Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		})
	}

	return products, nil
}
