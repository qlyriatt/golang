package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

func Update() {
	// запрос всех подов
	pods, err := g_deploy.GetPodList()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	// перевод массива в map для ускорения поиска по client_id
	m := make(map[string]struct{})
	for _, name := range pods {
		m[name] = struct{}{}
	}

	// запрос всех клиентов из БД
	rows, err := g_db.Query(context.Background(), "SELECT * FROM clients")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	// проход полученных из БД данных для выявления несовпадений статусов алгоритмов
	// в подах и в БД
	for rows.Next() {
		var client_id int64
		var vwap bool
		var twap bool
		var hft bool
		rows.Scan(nil, &client_id, &vwap, &twap, &hft)

		client_id_str := strconv.FormatInt(client_id, 10)

		_, ok := m[client_id_str+"_vwap"]
		if !ok && vwap {
			err := g_deploy.CreatePod(client_id_str + "_vwap")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		} else if ok && !vwap {
			err := g_deploy.DeletePod(client_id_str + "_vwap")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}

		_, ok = m[client_id_str+"_twap"]
		if !ok && twap {
			err := g_deploy.CreatePod(client_id_str + "_twap")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		} else if ok && !twap {
			err := g_deploy.DeletePod(client_id_str + "_twap")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}

		_, ok = m[client_id_str+"_hft"]
		if !ok && hft {
			err := g_deploy.CreatePod(client_id_str + "_hft")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		} else if ok && !hft {
			err := g_deploy.DeletePod(client_id_str + "_hft")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	}
}
