import Operation from "./Operation";
import { map } from 'lodash';
import Api from "../core/Api";
import { observable } from 'mobx';

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

  public loadOperation(id: number): Promise<Operation> {
    return Api.GET( `/operations/${id}`, {} ).then( ( response ) => {
      return new Operation( response.data.data.Operation );
    } );
  }

  public loadOperations(): Promise<Operation[]> {
    return Api.GET( '/operations', {} ).then( ( response ) => {
      this.operations = map( response.data.data.Operations, ( o: any ) => {
        return new Operation( o );
      } );

      return this.operations;
    } );
  }

  public saveOperation(operation: Operation): Promise<void> {
    return Api.POST(`/operations/${operation.ID}`, operation.getData());
  }

  public createOperation(operation: Operation): Promise<void> {
    return Api.POST(`/operations`, operation.getData());
  }

  public deleteOperation(operation: Operation): Promise<void> {
    return Api.POST(`/operations/${operation.ID}/delete`, {});
  }
}