package api

import (
	"encoding/base64"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UDPRelayAPIFormStruct struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Data string `json:"data"`
}

type UDPResponse struct {
	Data string `json:"data"`
}

func HandlePostUdpRelay(ctx *gin.Context) {
	apiForm := UDPRelayAPIFormStruct{}

	if err := ctx.ShouldBindBodyWithJSON(&apiForm); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	decodedData, err := base64.StdEncoding.DecodeString(apiForm.Data)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid base64 data")
		return
	}

	addr := net.JoinHostPort(apiForm.Host, strconv.Itoa(apiForm.Port))
	conn, err := net.Dial("udp", addr)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to dial UDP: "+err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write(decodedData)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to send UDP packet: "+err.Error())
		return
	}

	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		ctx.String(http.StatusGatewayTimeout, "Timeout or read error: "+err.Error())
		return
	}

	response := UDPResponse{
		Data: base64.StdEncoding.EncodeToString(buffer[:n]),
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, response)
}
