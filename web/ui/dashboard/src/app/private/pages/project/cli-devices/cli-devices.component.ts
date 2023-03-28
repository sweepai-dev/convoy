import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PrivateService } from 'src/app/private/private.service';
import { PAGINATION } from 'src/app/models/global.model';
import { DEVICE } from 'src/app/models/endpoint.model';

@Component({
	selector: 'convoy-cli-devices',
	standalone: true,
	imports: [CommonModule],
	templateUrl: './cli-devices.component.html',
	styleUrls: ['./cli-devices.component.scss']
})
export class CliDevicesComponent implements OnInit {
	isLoadindingdevices = false;
	devices!: { pagination: PAGINATION; content: DEVICE[] };

	constructor(public privateService: PrivateService) {}

	ngOnInit(): void {
		this.getDevices();
	}

	async getDevices() {
		this.isLoadindingdevices = true;
		// try {
		// 	const response = await this.deviceService.getAppDevices(this.endpointId, this.token);
		// 	this.devices = response.data.content;
		// 	this.isLoadindingdevices = false;
		// } catch {
		// 	this.isLoadindingdevices = false;
		// 	return;
		// }
	}
}
