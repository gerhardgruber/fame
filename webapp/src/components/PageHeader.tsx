import * as React from 'react';
import moment from 'moment';
import { Layout } from 'antd';
import { observer } from 'mobx-react';
import UiStore from '../stores/UiStore';

const { Header } = Layout;

export interface PageHeaderProps {
    name: string;
}

@observer
class PageHeader extends React.Component<PageHeaderProps, any> {
    private uiStore: UiStore;

    constructor( props ) {
        super( props );

        this.uiStore = UiStore.getInstance();
    }

    render() {
        return (
            <Header>
                <div style={{ float: "left" }}>
                    <h1 style={{ color: "rgba(255, 255, 255, 0.85)" }}>{this.uiStore.T("HEADING_" + this.props.name)}</h1>
                </div>
                {/* <div style={{ float: "right" }}>
                    <h4 style={{ color: "rgba(255, 255, 255, 0.85)" }}>{moment().format("HH:mm:ss DD.MM.Y")}</h4>
                </div> */}
            </Header>
        )
    }
}

export default PageHeader;