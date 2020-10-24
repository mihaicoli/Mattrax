// TODO: Pagination, Filters

export const actions = {
  getAll(_context: any) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/0a624a65-eb2c-4ca5-aeb2-822e6397dfd0')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching devices from server'))
            return
          }

          const devices = await res.json()
          resolve(devices)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(_context: any, _deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/45a4d87a-92fc-4f30-9ca3-5ef2bf284385')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching device from server'))
            return
          }

          const device = await res.json()
          resolve(device)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getInformationByID(_context: any, _deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/719e7183-bac2-4a5b-975c-5a2c51458d9d')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching device from server'))
            return
          }

          const deviceInfo = await res.json()
          resolve(deviceInfo)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getScopeByID(_context: any, _deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/1ae5b4e1-2a4a-49ea-aee4-f97274a1a371')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching device from server'))
            return
          }

          const deviceScope = await res.json()
          resolve(deviceScope)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
