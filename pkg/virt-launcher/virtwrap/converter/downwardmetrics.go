package converter

import (
	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/downwardmetrics"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

func convertDownwardMetricsChannel() api.Channel {
	return api.Channel{
		Type: "unix",
		Source: &api.ChannelSource{
			Mode: "bind",
			Path: downwardmetrics.DownwardMetricsChannelSocket,
		},
		Target: &api.ChannelTarget{
			Type: v1.VirtIO,
			Name: downwardmetrics.DownwardMetricsSerialDeviceName,
		},
	}
}
