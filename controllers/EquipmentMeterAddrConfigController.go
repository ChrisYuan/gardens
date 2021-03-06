package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/yunnet/gardens/enums"
	"github.com/yunnet/gardens/models"
	"strconv"
	"strings"
	"time"
)

type EquipmentMeterAddrConfigController struct {
	BaseController
}

func (this *EquipmentMeterAddrConfigController) Prepare() {
	this.BaseController.Prepare()
	this.checkAuthor("DataGrid", "DataList")
}

func (this *EquipmentMeterAddrConfigController) Index() {
	this.Data["pageTitle"] = "寄存器分段配置"
	this.Data["showMoreQuery"] = true

	this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	this.setTpl()
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["headcssjs"] = "equipmentmeteraddrconfig/index_headcssjs.html"
	this.LayoutSections["footerjs"] = "equipmentmeteraddrconfig/index_footerjs.html"

	//页面里按钮权限控制
	this.Data["canEdit"] = this.checkActionAuthor("EquipmentMeterAddrConfigController", "Edit")
	this.Data["canDelete"] = this.checkActionAuthor("EquipmentMeterAddrConfigController", "Delete")
}

func (this *EquipmentMeterAddrConfigController) DataGrid() {
	var params models.EquipmentMeterAddrConfigQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	data, total := models.EquipmentMeterAddrConfigPageList(&params)

	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *EquipmentMeterAddrConfigController) DataList() {
	var params = models.EquipmentMeterAddrConfigQueryParam{}
	data := models.EquipmentMeterAddrConfigDataList(&params)
	this.jsonResult(enums.JRCodeSucc, "", data)
}

func (this *EquipmentMeterAddrConfigController) Edit() {
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}

	Id, _ := this.GetInt(":id", 0)
	m := models.EquipmentMeterAddrConfig{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
	} else {
		m.Used = enums.Enabled
		m.SegmentStartAddr = 1
		m.SegmentNO = 1
	}

	this.Data["m"] = m
	this.setTpl("equipmentmeteraddrconfig/edit.html", "shared/layout_pullbox.html")
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["footerjs"] = "equipmentmeteraddrconfig/edit_footerjs.html"
}

//add | update
func (this *EquipmentMeterAddrConfigController) Save() {
	var err error
	m := models.EquipmentMeterAddrConfig{}

	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	//tmpInt := this.Input().Get("Id")
	//m.Id, _ = strconv.Atoi(tmpInt)
	//
	//m.MeterTypeNO = this.GetString("MeterTypeNO")
	//
	//tmpInt = this.Input().Get("SegmentStartAddr")
	//m.SegmentStartAddr, _ = strconv.Atoi(tmpInt)
	//
	//tmpInt = this.Input().Get("SegmentLen")
	//m.SegmentLen, _ = strconv.Atoi(tmpInt)
	//
	//tmpInt = this.Input().Get("SegmentNO")
	//m.SegmentNO, _ = strconv.Atoi(tmpInt)

	m.ChangeUser = this.curUser.RealName
	//m.ChangeDate = time.Now()

	o := orm.NewOrm()
	if m.Id == 0 {
		m.CreateUser = this.curUser.RealName
		m.CreateDate = time.Now()

		if _, err = o.Insert(&m); err == nil {
			this.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if _, err = o.Update(&m, "MeterTypeNO", "SegmentStartAddr", "SegmentLen", "SegmentNO", "Used", "ChangeUser", "ChangeDate"); err == nil {
			this.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			this.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

func (this *EquipmentMeterAddrConfigController) Delete() {
	strs := this.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.EquipmentMeterAddrConfigBatchDelete(ids); err == nil {
		this.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		this.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
