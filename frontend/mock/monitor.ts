import { defineFakeRoute } from "vite-plugin-fake-server/client";

export default defineFakeRoute([
  {
    url: "/api/monitor/server",
    method: "get",
    response: () => {
      return {
        code: 200,
        message: "success",
        payload: {
          general: {
            hostname: "DESKTOP-V87B08V",
            os: "Microsoft Windows 11 Enterprise LTSC 2024 10.0.26100.3775 Build 26100.3775",
            arch: "amd64",
            kernel: "10.0.26100.3775 Build 26100.3775",
            uptime: 12522514,
            uptimeFormat: "144天 22小时 28分钟",
            cpuCores: 8,
            goVersion: "go1.24.2"
          },
          cpu: {
            usagePercent: 30.859375,
            cores: 8,
            modelName: "Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz",
            frequency: 2112,
            coreDetails: [
              23.076923076923077, 16.666666666666664, 20.3125, 23.4375, 25,
              17.1875, 26.5625, 13.846153846153847
            ]
          },
          memory: {
            total: 17004486656,
            used: 12277092352,
            free: 4727394304,
            available: 4727394304,
            usedPercent: 72,
            totalFormat: "15.8 GB",
            usedFormat: "11.4 GB",
            freeFormat: "4.4 GB"
          },
          disk: [
            {
              device: "C:",
              mountpoint: "C:",
              fstype: "NTFS",
              total: 511037476864,
              used: 368528269312,
              free: 142509207552,
              usedPercent: 72.11374625076171,
              totalFormat: "475.9 GB",
              usedFormat: "343.2 GB",
              freeFormat: "132.7 GB"
            }
          ],
          network: [
            {
              name: "WLAN",
              bytesRecv: 6131787545,
              bytesSent: 10537761665,
              packetsRecv: 10980979,
              packetsSent: 15742714,
              recvFormat: "5.7 GB",
              sentFormat: "9.8 GB",
              isUp: true
            },
            {
              name: "Meta",
              bytesRecv: 420369203,
              bytesSent: 205669221,
              packetsRecv: 207999,
              packetsSent: 358665,
              recvFormat: "400.9 MB",
              sentFormat: "196.1 MB",
              isUp: true
            }
          ],
          load: {
            load1: 0.012907845951580859,
            load5: 0.4124272315231688,
            load15: 0.9078386720538169
          },
          process: {
            total: 19,
            running: 1,
            sleeping: 0,
            stopped: 0,
            zombie: 0
          },
          timestamp: 1757853494
        }
      };
    }
  }
]);
