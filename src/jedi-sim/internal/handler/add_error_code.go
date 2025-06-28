package handler

import (
	"encoding/json"
	"io/ioutil"
	"jedi-sim/internal/model"
	"log"
	"net/http"
	"os"
	"strconv"
)

// AddErrorCodeHandler 处理添加错误码数据的保存
func AddErrorCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[AddErrorCodeHandler] Request received:", r.Method, r.URL.Path)

	// 只接受POST请求
	if r.Method != "POST" {
		log.Println("[AddErrorCodeHandler] Invalid request method:", r.Method)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// 解析请求体中的JSON数据
	var errorCodeData model.ErrorInfo
	err := json.NewDecoder(r.Body).Decode(&errorCodeData)
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Failed to decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid JSON data",
		})
		return
	}

	// 读取现有的errorCodes.json文件
	filePath := "errorCodes.json"
	var existingData map[string]model.ErrorInfo

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果文件不存在，创建一个新的空map
		log.Printf("[AddErrorCodeHandler] Error reading file: %v, creating new data", err)
		existingData = make(map[string]model.ErrorInfo)
	} else {
		// 解析现有的JSON数据
		err = json.Unmarshal(fileData, &existingData)
		if err != nil {
			log.Printf("[AddErrorCodeHandler] Error unmarshaling JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  false,
				"message": "Error processing existing data",
			})
			return
		}
	}

	// 使用Z4Z5_ErrorCode作为键
	key := strconv.Itoa(errorCodeData.Z4Z5_ErrorCode)
	// 将新数据添加到现有数据中
	existingData[key] = errorCodeData

	// 将更新后的数据写回文件
	updatedData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Error processing data",
		})
		return
	}

	err = ioutil.WriteFile(filePath, updatedData, os.ModePerm)
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Error writing file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Error saving data",
		})
		return
	}

	// 返回成功响应
	log.Println("[AddErrorCodeHandler] Data saved successfully")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Error code data saved successfully",
	})
}