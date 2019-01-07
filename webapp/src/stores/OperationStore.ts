import { map } from 'lodash';
import Api from "../core/Api";
import { observable } from 'mobx';
import Operation from "antd/lib/transfer/operation";

export default class OperationStore {
  private static instance: OperationStore;

  @observable public operations: Operation[];

  constructor() {
    OperationStore.instance = this;
  }

  public static getInstance(): OperationStore {
    if ( OperationStore.instance ) {
      return OperationStore.instance;
    }

    OperationStore.instance = new OperationStore();
    return OperationStore.instance;
  }

  public loadOperations(): Promise<Operation[]> {
    return Api.GET( '/operations', {} ).then( ( response ) => {
      this.operations = map( response.data.data.Users, ( o: any ) => {
        return new Operation( o );
      } );

      return this.operations;
    } );
  }
}