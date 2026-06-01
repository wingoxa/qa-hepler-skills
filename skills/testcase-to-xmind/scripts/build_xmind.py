#!/usr/bin/env python3
"""Build an old-format XMind file from structured test case JSON."""

from __future__ import annotations

import argparse
import json
import time
import uuid
import zipfile
from pathlib import Path
from typing import Any
from xml.etree import ElementTree as ET


CONTENT_NS = "urn:xmind:xmap:xmlns:content:2.0"
ET.register_namespace("", CONTENT_NS)
ET.register_namespace("fo", "http://www.w3.org/1999/XSL/Format")
ET.register_namespace("svg", "http://www.w3.org/2000/svg")
ET.register_namespace("xhtml", "http://www.w3.org/1999/xhtml")
ET.register_namespace("xlink", "http://www.w3.org/1999/xlink")


def new_id() -> str:
    return uuid.uuid4().hex[:26]


def now_ms() -> str:
    return str(int(time.time() * 1000))


def text(value: Any, default: str = "") -> str:
    if value is None:
        return default
    if isinstance(value, str):
        return value.strip()
    return str(value).strip()


def normalize_priority(value: Any) -> str:
    priority = text(value, "P2").upper()
    return priority if priority in {"P0", "P1", "P2", "P3"} else "P2"


def normalize_tags(value: Any) -> str:
    if value is None:
        return ""
    if isinstance(value, list):
        return "｜".join(text(item) for item in value if text(item))
    return text(value)


def sub(parent: ET.Element, name: str, attrs: dict[str, str] | None = None) -> ET.Element:
    return ET.SubElement(parent, name, attrs or {})


def topic(title: str, children: list[ET.Element] | None = None, style_id: str = "42gut2jr74kvarqfndgjugps26") -> ET.Element:
    node = ET.Element(
        "topic",
        {
            "id": new_id(),
            "style-id": style_id,
            "timestamp": now_ms(),
        },
    )
    title_node = sub(node, "title")
    title_node.text = title
    if children:
        children_node = sub(node, "children")
        topics_node = sub(children_node, "topics", {"type": "attached"})
        for child in children:
            topics_node.append(child)
    return node


def case_topic(case: dict[str, Any]) -> ET.Element:
    name = text(case.get("name"), "未命名用例")
    children: list[ET.Element] = []

    precondition = text(case.get("precondition"))
    if precondition:
        children.append(topic(f"前置条件：{precondition}"))

    remark = text(case.get("remark"))
    if remark:
        children.append(topic(f"备注：{remark}"))

    tags = normalize_tags(case.get("tags"))
    if tags:
        children.append(topic(f"标签：{tags}"))

    steps = case.get("steps") or []
    if isinstance(steps, list) and steps:
        step_nodes: list[ET.Element] = []
        for item in steps:
            if not isinstance(item, dict):
                item = {"action": item}
            action = text(item.get("action") or item.get("step"), "待补充")
            expected = text(item.get("expected") or item.get("result"))
            expected_nodes = [topic(f"预期结果：{expected}")] if expected else []
            step_nodes.append(topic(f"步骤：{action}", expected_nodes))
        children.append(topic("步骤描述", step_nodes))
    else:
        description = text(case.get("description") or case.get("step") or case.get("steps"))
        expected = text(case.get("expected") or case.get("result"))
        if description or expected:
            expected_nodes = [topic(f"预期结果：{expected}")] if expected else []
            children.append(topic(f"文本描述：{description or '待补充'}", expected_nodes))

    children.append(topic(f"用例等级：{normalize_priority(case.get('priority'))}"))
    return topic(f"case：{name}", children)


def build_content(data: dict[str, Any]) -> bytes:
    root = ET.Element(
        f"{{{CONTENT_NS}}}xmap-content",
        {
            "timestamp": now_ms(),
            "version": "2.0",
        },
    )
    sheet = sub(
        root,
        "sheet",
        {
            "xmlns": "",
            "id": new_id(),
            "style-id": "4k148f1i7ntd91ftp26ks67r4u",
            "timestamp": now_ms(),
        },
    )

    modules = data.get("modules") or []
    module_nodes: list[ET.Element] = []
    for module in modules:
        if not isinstance(module, dict):
            continue
        module_name = text(module.get("name"), "默认模块")
        case_nodes = [case_topic(case) for case in (module.get("cases") or []) if isinstance(case, dict)]
        module_nodes.append(topic(module_name, case_nodes, "5239lu9r6doet0mjj3438h5a41"))

    root_topic = topic(text(data.get("title"), "功能用例"), module_nodes, "1529ujvcb63rhopv6kjv4pg7dt")
    sheet.append(root_topic)
    return ET.tostring(root, encoding="utf-8", xml_declaration=True)


def build_xmind(input_path: Path, template_path: Path, output_path: Path) -> None:
    data = json.loads(input_path.read_text(encoding="utf-8"))
    content_xml = build_content(data)

    output_path.parent.mkdir(parents=True, exist_ok=True)
    with zipfile.ZipFile(template_path, "r") as zin:
        with zipfile.ZipFile(output_path, "w", zipfile.ZIP_DEFLATED) as zout:
            replaced = False
            for item in zin.infolist():
                if item.filename == "content.xml":
                    zout.writestr(item, content_xml)
                    replaced = True
                else:
                    zout.writestr(item, zin.read(item.filename))
            if not replaced:
                zout.writestr("content.xml", content_xml)


def main() -> None:
    skill_dir = Path(__file__).resolve().parents[1]
    parser = argparse.ArgumentParser(description="Build an XMind test case map from JSON.")
    parser.add_argument("--input", required=True, type=Path, help="Path to test case JSON.")
    parser.add_argument("--output", required=True, type=Path, help="Output .xmind path.")
    parser.add_argument(
        "--template",
        type=Path,
        default=skill_dir / "assets" / "xmind_case.xmind",
        help="Template .xmind path.",
    )
    args = parser.parse_args()

    build_xmind(args.input, args.template, args.output)
    print(f"Generated {args.output}")


if __name__ == "__main__":
    main()
