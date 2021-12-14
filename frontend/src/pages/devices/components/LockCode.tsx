import React from 'react';
import { Button, Form } from 'react-bootstrap';
import { addDays, format, parseISO, set } from 'date-fns';
import { Loading } from '../../utils/Loading';
import { StandardFetch } from '../../utils/FetchHelper';

interface Props{
    deviceId: string;
    managedLockCode: DeviceManagedLockCodeT | null;
    managedLockCodesUpdated: () => void;
}

export const LockCode = (props:Props) => {
    const [loading, setLoading] = React.useState<boolean>(false);
    const [code, setCode] = React.useState<string>(
        props.managedLockCode ? props.managedLockCode.code + " - " + props.managedLockCode.status + " - " + props.managedLockCode.note : ""
    );
    const [startAt, setStartAt] = React.useState<Date>(props.managedLockCode ? parseISO(props.managedLockCode.startAt) : new Date());
    const [endAt, setEndAt] = React.useState<Date>(props.managedLockCode ? parseISO(props.managedLockCode.endAt) : addDays(set(new Date(), {minutes: 0}), 1));

    const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault();
        setLoading(true);

        if (props.managedLockCode === null) {
            StandardFetch("devices/" + props.deviceId + "/lock-codes/", {
                method: "POST",
                body: JSON.stringify({
                    code: code,
                    startAt: startAt,
                    endAt: endAt,
                })
            })
            .then(response => response.json())
            .then(response => {
                setCode("");
                props.managedLockCodesUpdated();
                setLoading(false);
            })
            .catch(err => {
                // TODO: indicate error.
                console.log(err);
            });
            return;
        }

        StandardFetch("devices/" + props.deviceId + "/lock-codes/" + props.managedLockCode?.id, {
            method: "PUT",
            body: JSON.stringify({
                endAt: endAt,
            })
        })
        .then(response => response.json())
        .then(response => {
            props.managedLockCodesUpdated();
            setLoading(false);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    };

    const render = () => {
        // Might have just a "renderLockCode" method in the future too?
        return renderAddLockCode();
    };

    const renderAddLockCode = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <Form onSubmit={evt => formSubmit(evt)}>

                <Form.Group>
                    <Form.Label>Code</Form.Label>
                    <Form.Control type="text" defaultValue={code} onChange={(evt) => setCode(evt.target.value)} disabled={!!props.managedLockCode}/>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Enable At</Form.Label>
                    <Form.Control type="datetime-local" defaultValue={format(startAt, "yyyy-MM-dd'T'HH:mm")} onChange={(evt) => setStartAt(parseISO(evt.target.value))} disabled={!!props.managedLockCode}/>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Disable At</Form.Label>
                    <Form.Control type="datetime-local" defaultValue={format(endAt, "yyyy-MM-dd'T'HH:mm")} onChange={(evt) => setEndAt(parseISO(evt.target.value))}/>
                </Form.Group>

                <Button variant="secondary" type="submit">
                    {props.managedLockCode ? "Update" : "Add"} Lock Code
                </Button>

            </Form>
        );
    };

    return render();
};
