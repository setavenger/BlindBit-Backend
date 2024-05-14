package server

import (
	"SilentPaymentAppBackend/src/common"
	"SilentPaymentAppBackend/src/common/types"
	"SilentPaymentAppBackend/src/db/dblevel"
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// ApiHandler todo might not need ApiHandler struct if no data is stored within.
//
//	Will keep for now just in case, so I don't have to refactor twice
type ApiHandler struct{}

type TxRequest struct {
	Data string `form:"data" json:"data" binding:"required"`
}

func (h *ApiHandler) GetBestBlockHeight(c *gin.Context) {
	// todo returns one height too low
	lastHeader, err := dblevel.FetchHighestBlockHeaderInvByFlag(true)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could could not retrieve data from database",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"block_height": lastHeader.Height,
	})
}

func (h *ApiHandler) GetCFilterByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}

	cFilter, err := dblevel.FetchByBlockHashFilter(headerInv.Hash)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get filter from db",
		})
		return
	}

	data := gin.H{
		"filter_type":  cFilter.FilterType,
		"block_height": height, // saves us a "join" in the query
		"block_hash":   cFilter.BlockHash,
		"data":         hex.EncodeToString(cFilter.Data),
	}

	c.JSON(200, data)
}

func (h *ApiHandler) GetUtxosByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}
	utxos, err := dblevel.FetchByBlockHashUTXOs(headerInv.Hash)
	if err != nil && !errors.Is(err, dblevel.NoEntryErr{}) {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could could not retrieve data from database",
		})
		return
	}
	if utxos != nil {
		c.JSON(200, utxos)
	} else {
		c.JSON(200, []interface{}{})
	}
}

// GetTweakDataByHeight serves tweak data as json array of tweaks (33 byte as hex-formatted)
// todo can be changed to serve with verbosity aka serve with txid or even block data (height, hash)
func (h *ApiHandler) GetTweakDataByHeight(c *gin.Context) {
	// todo outsource all the blockHeight extraction and conversion through the inverse header table into middleware
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad format: height",
		})
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	// Extracting the dustLimit query parameter and converting it to uint64
	dustLimitStr := c.DefaultQuery("dustLimit", "0") // Default to "0" if not provided
	dustLimit, err := strconv.ParseUint(dustLimitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dustLimit parameter"})
		return
	}

	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}
	var tweaks []types.Tweak

	if dustLimit == 0 {
		// this query should have a better performance due to no required checks
		tweaks, err = dblevel.FetchByBlockHashTweaks(headerInv.Hash)
		if err != nil && !errors.Is(err, dblevel.NoEntryErr{}) {
			common.ErrorLogger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could could not retrieve data from database",
			})
			return
		}
	} else {
		tweaks, err = dblevel.FetchByBlockHashDustLimitTweaks(headerInv.Hash, dustLimit)
		if err != nil && !errors.Is(err, dblevel.NoEntryErr{}) {
			common.ErrorLogger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could could not retrieve data from database",
			})
			return
		}
	}

	if err != nil && errors.Is(err, dblevel.NoEntryErr{}) {
		c.JSON(http.StatusOK, []string{})
		return
	}

	var serveTweakData []string
	for _, tweak := range tweaks {
		serveTweakData = append(serveTweakData, hex.EncodeToString(tweak.TweakData[:]))
	}

	c.JSON(http.StatusOK, serveTweakData)
}

func (h *ApiHandler) GetTweakIndexDataByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad format: height",
		})
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}

	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}

	// Extracting the dustLimit query parameter and converting it to uint64
	dustLimitStr := c.DefaultQuery("dustLimit", "0") // Default to "0" if not provided
	dustLimit, err := strconv.ParseUint(dustLimitStr, 10, 64)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dustLimit parameter"})
		return
	}

	if dustLimit != 0 && !common.TweakIndexFullIncludingDust {
		common.DebugLogger.Println("tried accessing dust limits")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server does not allow dustLimits"})
		return
	}

	// todo basically duplicate code could be simplified and generalised with interface/(generics?)
	if common.TweakIndexFullIncludingDust {
		var tweakIndex *types.TweakIndexDust
		tweakIndex, err = dblevel.FetchByBlockHashTweakIndexDust(headerInv.Hash)
		if err != nil && !errors.Is(err, dblevel.NoEntryErr{}) {
			common.ErrorLogger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could could not retrieve data from database",
			})
			return
		}

		if err != nil && errors.Is(err, dblevel.NoEntryErr{}) {
			c.JSON(http.StatusOK, []string{})
			return
		}

		var serveTweakData = []string{}
		for _, tweak := range tweakIndex.Data {
			common.DebugLogger.Printf("%x- %d", tweak.Tweak(), tweak.HighestValue())
			if tweak.HighestValue() < dustLimit {
				continue
			}
			data := tweak.Tweak()
			serveTweakData = append(serveTweakData, hex.EncodeToString(data[:]))
		}

		c.JSON(200, serveTweakData)
		return
	} else {
		// this query should have a better performance due to no required checks
		tweakIndex, err := dblevel.FetchByBlockHashTweakIndex(headerInv.Hash)
		if err != nil && !errors.Is(err, dblevel.NoEntryErr{}) {
			common.ErrorLogger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could could not retrieve data from database",
			})
			return
		}

		if err != nil && errors.Is(err, dblevel.NoEntryErr{}) {
			c.JSON(http.StatusOK, []string{})
			return
		}

		var serveTweakData []string
		for _, tweak := range tweakIndex.Data {
			serveTweakData = append(serveTweakData, hex.EncodeToString(tweak[:]))
		}

		c.JSON(200, serveTweakData)
		return
	}
}

func (h *ApiHandler) ForwardRawTX(c *gin.Context) {
	var txRequest TxRequest
	if err := c.ShouldBind(&txRequest); err != nil {
		common.ErrorLogger.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	err := forwardTxToMemPool(txRequest.Data)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func forwardTxToMemPool(txHex string) error {
	//url := "http://localhost/api/tx"

	resp, err := http.Post(common.MempoolEndpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(txHex))
	if err != nil {
		common.ErrorLogger.Printf("Failed to make request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.ErrorLogger.Printf("Failed to read response: %s\n", err)
		return err
	}

	common.DebugLogger.Println("Response:", string(body))
	return nil
}
