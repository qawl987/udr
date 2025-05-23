/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package processor

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/internal/util"
	"github.com/free5gc/util/mongoapi"
)

func (p *Processor) QuerySmDataProcedure(c *gin.Context, collName string, ueId string, servingPlmnId string,
	singleNssai models.Snssai, dnn string,
) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	if !reflect.DeepEqual(singleNssai, models.Snssai{}) {
		if singleNssai.Sd == "" {
			filter["singleNssai.sst"] = singleNssai.Sst
		} else {
			filter["singleNssai.sst"] = singleNssai.Sst
			filter["singleNssai.sd"] = singleNssai.Sd
		}
	}

	if dnn != "" {
		dnnKey := util.EscapeDnn(dnn)
		filter["dnnConfigurations."+dnnKey] = bson.M{"$exists": true}
	}
	resp := models.SmSubsData{}

	sessionManagementSubscriptionDatas, err := mongoapi.
		RestfulAPIGetMany(collName, filter, mongoapi.COLLATION_STRENGTH_IGNORE_CASE)
	if err != nil {
		logger.DataRepoLog.Errorf("QuerySmDataProcedure err: %+v", err)
		pd := util.ProblemDetailsUpspecified("")
		c.JSON(int(pd.Status), pd)
		return
	}
	for _, smData := range sessionManagementSubscriptionDatas {
		var tmpSmData models.SessionManagementSubscriptionData
		err := json.Unmarshal(util.MapToByte(smData), &tmpSmData)
		if err != nil {
			logger.DataRepoLog.Debug("SmData Unmarshal error")
			continue
		}
		resp.IndividualSmSubsData = append(resp.IndividualSmSubsData, tmpSmData)

		dnnConfigurations := tmpSmData.DnnConfigurations
		tmpDnnConfigurations := make(map[string]models.DnnConfiguration)
		for escapedDnn, dnnConf := range dnnConfigurations {
			dnn := util.UnescapeDnn(escapedDnn)
			tmpDnnConfigurations[dnn] = dnnConf
		}
		smData["DnnConfigurations"] = tmpDnnConfigurations
	}
	c.JSON(http.StatusOK, resp)
}
