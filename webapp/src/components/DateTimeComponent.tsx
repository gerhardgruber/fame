import * as React from 'react';
import Moment from 'react-moment';
import UiStore from '../stores/UiStore';

interface DateTimeComponentProps {
    datetime: any;
    format: string;
}

class DateTimeComponent extends React.Component<DateTimeComponentProps, never> {

    static makeComponent(datetime: any, format: string) {
        return <DateTimeComponent datetime={datetime} format={format} />;
    }

    render() {
        return <Moment date={this.props.datetime} format={this.props.format} />;
    }
}

export default DateTimeComponent;