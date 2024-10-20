import wretch from 'wretch'
import QueryStringAddon from 'wretch/addons/queryString'

const url = '/api'

export const http = wretch(url).addon(QueryStringAddon)
