import { API_ROOT } from '../../config/fame'
import axios from 'axios'
import isObject from 'lodash/isObject'
import isString from 'lodash/isString'
import isFunction from 'lodash/isFunction'
import isArray from 'lodash/isArray'
import isNil from 'lodash/isNil'
import each from 'lodash/each'
import find from 'lodash/find'
import has from 'lodash/has'
import template from 'lodash/template'
import { toJS } from 'mobx'

import StateStore from '../stores/StateStore'
import UiStore from '../stores/UiStore'

const state = StateStore.getInstance();

namespace Api
{
  export const callback = {
    translate: (caption: string, data: any = null): string => {
      return caption;
    }
  }

  function makeParams( method : string, url : string, opts ) {
    const params = { method, url, baseURL: API_ROOT }

    const payload = { ...opts,
      responseType: 'json'
    }

    if( method === 'POST' || method === 'PUT' || method === 'PATCH' ) {
      let fd = payload

      if( hasInstanceOf( File, opts ) ) {
        fd = new FormData()
        each( payload, ( v, k ) => {
          if( isObject( v ) && !( v instanceof File) ) {
            fd.append( k, JSON.stringify( v ) )
          } else {
            fd.append( k, v )
          }
        } )
      }

      params[ 'data' ] = fd
    } else {
      params[ 'params' ] = payload
    }

    params[ 'params' ] = { ...params[ 'params' ],
      session: window.sessionStorage.getItem( 'session' )
    }

    return params
  }

  export function HEAD( endpoint: string, args: object ) : Promise<any> {
    return call( 'HEAD', endpoint, args );
  }

  export function GET( endpoint: string, args: object ) : Promise<any> {
    return call( 'GET', endpoint, args );
  }

  export function POST( endpoint: string, args: object ) : Promise<any> {
    return call( 'POST', endpoint, args );
  }

  export function PUT( endpoint: string, args: object ) : Promise<any> {
    return call( 'PUT', endpoint, args );
  }

  export function PATCH( endpoint: string, args: object ) : Promise<any> {
    return call( 'PATCH', endpoint, args );
  }

  export function DELETE( endpoint: string, args: object ) : Promise<any> {
    return call( 'DELETE', endpoint, args );
  }

  function call( method : string, endpoint: string, args: object ) : Promise<any> {
    beforeApiStart()

    const params = makeParams( method, endpoint, args );

    return request( params )
  }

  function request( params ) : Promise<any> {
    return axios.request( params ).catch( systemError ).then( successHandler );
  }

  function successHandler( response ) {
    state.progressAfterLoading()
    handleSuccess( response )
    afterApiEnd()

    return response;
  }

  /**
   * Performing multiple concurrent requests, waiting until last response arrived.
   * @param  {Array}   Promises An array of requests (Promises)
   * @return {Promise} Returns Promise after all requests have finished.
   */
  export function whenAll( ...promises ) {
    return axios.all( promises )
  }

  function hasInstanceOf( Class, objects ) {
    let matchFound = false

    find( objects, obj => {
      if ( obj instanceof Class ) {
        return matchFound = true
      }
    })

    return matchFound
  }


  function beforeApiStart() {
    state.startLoading()
  }

  function afterApiEnd() {
    state.endLoading()
  }

  function handleSuccess( response ) {
    if( !has( response, 'data.success' ) ) { return }

    switch( response.data.success ) {
    case 'session':
      handleExpired( response.data );
      throw "Session expired!";

    case true:
      break

    case false:
      handleError( response.data );
    }

    return response
  }

  function handleExpired( error ) {
    // TODO: extract to a session class
    window.sessionStorage.setItem( 'session', null )
    afterApiEnd()
  }

  function handleError( error ) {
    switch( error.message ) {
    case "AuthenticationError":
      handleExpired( error )
    default:
      state.notifications.error( error.message )
    }
    throw error.message || "Request error!";
  }

  function systemError( error ) {
    state.notifications.error( callback.translate(
      error.response.data.caption,
      error.response.data.captionData
     ) )
    afterApiEnd()
    throw error
  }

  function notify( data ) {
    if( !isString( data.msg ) ) {
      return
    }

    if( data.success === true ) {
      state.notifications.info( data.msg )
    } else {
      state.notifications.error( data.msg )
    }
  }
}

export default Api;