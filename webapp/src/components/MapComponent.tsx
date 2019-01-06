import * as React from 'react';
import { compose, withProps } from 'recompose';
import { withScriptjs, withGoogleMap, GoogleMap } from 'react-google-maps'

const DEFAULT_ZOOM = 8;
const MAPS_API_KEY = "AIzaSyDy6CKfFfgtFkgnltRkvvJy5F9ttsAHs60";
const BASE_URL = "https://maps.googleapis.com/maps/api/js?v=3.exp&key=" + MAPS_API_KEY;

interface MapComponentProps {
    // Default 8
    defaultZoom?: number;
    // If not given and takes the center of all child elements with position
    defaultCenter?: { lat: number, lng: number };

    // If not given takes a div that fills the container element
    mapElement?: JSX.Element;
    // If not given takes a div that fills the container element
    loadingElement?: JSX.Element;
    // If not given takes a div with height 400px
    containerElement?: JSX.Element;

    // If not given is ["geometry","drawing","places"]
    libraries?: Array<string>;
}

class _InnerMapComponent extends React.Component<MapComponentProps, any> {
    render() {
        return (
            <GoogleMap
                defaultZoom={this.props.defaultZoom}
                defaultCenter={this.props.defaultCenter}>

                {this.props.children}
            </GoogleMap>
        );
    }
}

const InnerMapComponent = withScriptjs(withGoogleMap(_InnerMapComponent));

class MapComponent extends React.Component<MapComponentProps, any> {
    private defaultZoom: MapComponentProps["defaultZoom"] = DEFAULT_ZOOM;
    private defaultCenter: MapComponentProps["defaultCenter"];

    private mapElement: MapComponentProps["mapElement"] = <div style={{ height: `100%` }} />;
    private loadingElement: MapComponentProps["loadingElement"] = <div style={{ height: `100%` }} />;
    private containerElement: MapComponentProps["containerElement"] = <div style={{ height: `400px` }} />;

    private libraries: MapComponentProps["libraries"] = ["geometry", "drawing", "places"];

    constructor(props: MapComponentProps) {
        super(props);
        this.componentWillReceiveProps(props);
    }

    componentWillReceiveProps(props) {
        this.defaultZoom = props.defaultZoom || this.defaultZoom;

        this.mapElement = props.mapElement || this.mapElement;
        this.loadingElement = props.loadingElement || this.loadingElement;
        this.containerElement = props.containerElement || this.containerElement;

        this.libraries = props.libraries || this.libraries;

        if (props.defaultCenter) {
            this.defaultCenter = props.defaultCenter;
        } else if (React.Children) {
            const children = React.Children.toArray(this.props.children);
            const avgPosition = { lat: 0, lng: 0 };
            let count = 0;

            children.forEach(child => {

                if (!this.defaultCenter) {

                    const childObject = (child as object);
                    if (childObject.hasOwnProperty('props') && childObject['props'].hasOwnProperty('position')
                        && childObject['props']['position'].hasOwnProperty('lat') && childObject['props']['position'].hasOwnProperty('lng')) {

                        avgPosition.lat += childObject['props']['position']['lat'];
                        avgPosition.lng += childObject['props']['position']['lng'];
                        count++;
                    }
                }
            });

            if (count > 0) {
                avgPosition.lat /= count;
                avgPosition.lng /= count;

                this.defaultCenter = avgPosition;
            }
        }

        if (!this.defaultCenter) {
            console.error("ERROR: No defaultCenter or child with position given to MapComponent");
        }
    }

    render() {
        return (
            <InnerMapComponent
                defaultZoom={this.defaultZoom}
                defaultCenter={this.defaultCenter}

                mapElement={this.mapElement}
                loadingElement={this.loadingElement}
                containerElement={this.containerElement}
                googleMapURL={BASE_URL + "&libraries=" + this.libraries.join()}>

                {this.props.children}
            </InnerMapComponent>
        );
    }
}

export default MapComponent;