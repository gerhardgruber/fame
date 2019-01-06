import Api from '../core/Api';
import Position from './Position';

class PositionStore {
    private static instance: PositionStore;

    public static getInstance() : PositionStore {
        if ( PositionStore.instance == null ) {
            PositionStore.instance = new PositionStore();
        }

        return PositionStore.instance;
    }
}

export default PositionStore;