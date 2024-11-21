import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get<{}, {}, {}, '/api/processors'>('/api/processors', async () => {
    return HttpResponse.json({
      'filter': {
        'endpoint_key': 'filter',
        'name': 'Filter',
        'description': 'Filters events based on field values',
        'parameters': {
          'inverse': {
            'description': 'Whether to invert the filter',
            'type': 1,
            'default_value': 'false'
          },
          'selector': {
            'description': 'The value to search for',
            'type': 3,
            'default_value': ''
          },
          'selectorFields': {
            'description': 'Whitespace separated list of event fields to filter on',
            'type': 4,
            'default_value': 'SUMMARY DESCRIPTION'
          }
        }
      }
    })
  })
]
