// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/policies', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching policies from server'))
            return
          }

          const policies = await res.json()
          resolve(policies)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, policyID: string) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/policy/' + encodeURI(policyID), {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching policy from server'))
            return
          }

          const policy = await res.json()
          resolve(policy)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
